package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_jobs_GetStatus(t *testing.T) {
	type serverResponse struct {
		Response   jobsResponse
		StatusCode int
	}

	tests := []struct {
		name     string
		storeID  string
		job      jobType
		response serverResponse
		want     jobStatus
		wantErr  bool
	}{
		{
			name:    "status done",
			storeID: "123",
			job:     deployStoreJob,
			response: serverResponse{
				Response: jobsResponse{
					Jobs: []jobsResponseJobItem{
						{
							JobID:  "456",
							Type:   deployStoreJob,
							Status: jobStatusFinished,
						},
						{
							JobID:  "789",
							Type:   modifyDbConfigJob,
							Status: jobStatusRunning,
						},
					},
					Total: 2,
				},
				StatusCode: http.StatusOK,
			},
			want:    jobStatusFinished,
			wantErr: false,
		},
		{
			name:    "status running",
			storeID: "123",
			job:     deployStoreJob,
			response: serverResponse{
				Response: jobsResponse{
					Jobs: []jobsResponseJobItem{
						{
							JobID:  "456",
							Type:   deployStoreJob,
							Status: jobStatusRunning,
						},
					},
					Total: 1,
				},
				StatusCode: http.StatusOK,
			},
			want:    jobStatusRunning,
			wantErr: false,
		},
		{
			name:    "job failed",
			storeID: "123",
			job:     deployStoreJob,
			response: serverResponse{
				Response: jobsResponse{
					Jobs: []jobsResponseJobItem{
						{
							JobID:  "456",
							Type:   deployStoreJob,
							Status: jobStatusErrored,
						},
					},
					Total: 1,
				},
				StatusCode: http.StatusOK,
			},
			want:    jobStatusErrored,
			wantErr: false,
		},
		{
			name:    "http error",
			storeID: "123",
			response: serverResponse{
				Response:   jobsResponse{},
				StatusCode: http.StatusInternalServerError,
			},
			want:    jobStatusUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const authToken = "fake-auth-token"
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/deployment/v2/data-stores/"+tt.storeID+"/jobs", r.URL.Path)
				require.Equal(t, authToken, r.Header.Get("Authorization"))

				w.WriteHeader(tt.response.StatusCode)
				err := json.NewEncoder(w).Encode(tt.response.Response)

				if err != nil {
					panic(err)
				}
			}))

			defer srv.Close()

			svc := jobs{
				baseURL: srv.URL,
				auth:    fakeAuthorizer{wantToken: authToken},
			}

			ctx := context.Background()

			got, err := svc.GetStatus(ctx, tt.storeID, tt.job)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jobs_Await(t *testing.T) {
	type serverResponse struct {
		Response   jobsResponse
		StatusCode int
	}

	tests := []struct {
		name      string
		storeID   string
		job       jobType
		responses []serverResponse
		want      jobStatus
		wantErr   bool
	}{
		{
			name:    "single time: job found, status done",
			storeID: "123",
			job:     deployStoreJob,
			responses: []serverResponse{
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   deployStoreJob,
								Status: jobStatusFinished,
							},
							{
								JobID:  "789",
								Type:   modifyDbConfigJob,
								Status: jobStatusRunning,
							},
						},
						Total: 2,
					},
					StatusCode: http.StatusOK,
				},
			},
			want:    jobStatusFinished,
			wantErr: false,
		},
		{
			name:    "job running, then job done",
			storeID: "123",
			job:     deployStoreJob,
			responses: []serverResponse{
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   deployStoreJob,
								Status: jobStatusRunning,
							},
						},
						Total: 1,
					},
					StatusCode: http.StatusOK,
				},
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   deployStoreJob,
								Status: jobStatusFinished,
							},
						},
						Total: 1,
					},
					StatusCode: http.StatusOK,
				},
			},
			want:    jobStatusFinished,
			wantErr: false,
		},
		{
			name:    "job running, then job failed",
			storeID: "123",
			job:     deployStoreJob,
			responses: []serverResponse{
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   deployStoreJob,
								Status: jobStatusRunning,
							},
						},
						Total: 1,
					},
					StatusCode: http.StatusOK,
				},
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   deployStoreJob,
								Status: jobStatusErrored,
							},
						},
						Total: 1,
					},
					StatusCode: http.StatusOK,
				},
			},
			want:    jobStatusErrored,
			wantErr: false,
		},
		{
			name:    "http error",
			storeID: "123",
			job:     deployStoreJob,
			responses: []serverResponse{
				{
					Response:   jobsResponse{},
					StatusCode: http.StatusGatewayTimeout,
				},
			},
			want:    jobStatusUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const authToken = "fake-auth-token"
			i := 0

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/deployment/v2/data-stores/"+tt.storeID+"/jobs", r.URL.Path)
				require.Equal(t, authToken, r.Header.Get("Authorization"))

				if len(tt.responses) == 0 {
					t.Fatalf("unexpected call to server")
				}

				if l := len(tt.responses); i >= l { // repeat the last response
					i = l - 1
				}

				w.WriteHeader(tt.responses[i].StatusCode)
				err := json.NewEncoder(w).Encode(tt.responses[i].Response)

				if err != nil {
					panic(err)
				}

				i++
			}))

			defer srv.Close()

			svc := jobs{
				baseURL:             srv.URL,
				auth:                fakeAuthorizer{wantToken: authToken},
				awaitTickerDuration: time.Second / 2,
			}

			ctx := context.Background()

			got, err := svc.Await(ctx, tt.storeID, tt.job, time.Second*10)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
