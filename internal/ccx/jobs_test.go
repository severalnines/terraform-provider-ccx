package ccx

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
		job      JobType
		response serverResponse
		want     JobStatus
		wantErr  bool
	}{
		{
			name:    "status done",
			storeID: "123",
			job:     DeployStoreJob,
			response: serverResponse{
				Response: jobsResponse{
					Jobs: []jobsResponseJobItem{
						{
							JobID:  "456",
							Type:   DeployStoreJob,
							Status: JobStatusFinished,
						},
						{
							JobID:  "789",
							Type:   ModifyDbConfigJob,
							Status: JobStatusRunning,
						},
					},
					Total: 2,
				},
				StatusCode: http.StatusOK,
			},
			want:    JobStatusFinished,
			wantErr: false,
		},
		{
			name:    "status running",
			storeID: "123",
			job:     DeployStoreJob,
			response: serverResponse{
				Response: jobsResponse{
					Jobs: []jobsResponseJobItem{
						{
							JobID:  "456",
							Type:   DeployStoreJob,
							Status: JobStatusRunning,
						},
					},
					Total: 1,
				},
				StatusCode: http.StatusOK,
			},
			want:    JobStatusRunning,
			wantErr: false,
		},
		{
			name:    "job failed",
			storeID: "123",
			job:     DeployStoreJob,
			response: serverResponse{
				Response: jobsResponse{
					Jobs: []jobsResponseJobItem{
						{
							JobID:  "456",
							Type:   DeployStoreJob,
							Status: JobStatusErrored,
						},
					},
					Total: 1,
				},
				StatusCode: http.StatusOK,
			},
			want:    JobStatusErrored,
			wantErr: false,
		},
		{
			name:    "http error",
			storeID: "123",
			response: serverResponse{
				Response:   jobsResponse{},
				StatusCode: http.StatusInternalServerError,
			},
			want:    JobStatusUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/deployment/v2/data-stores/"+tt.storeID+"/jobs", r.URL.Path)

				w.WriteHeader(tt.response.StatusCode)
				err := json.NewEncoder(w).Encode(tt.response.Response)

				if err != nil {
					panic(err)
				}
			}))

			defer srv.Close()

			svc := JobsClient{
				httpcli: NewTestHTTPClient(srv.URL),
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
		job       JobType
		responses []serverResponse
		want      JobStatus
		wantErr   bool
	}{
		{
			name:    "single time: job found, status done",
			storeID: "123",
			job:     DeployStoreJob,
			responses: []serverResponse{
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   DeployStoreJob,
								Status: JobStatusFinished,
							},
							{
								JobID:  "789",
								Type:   ModifyDbConfigJob,
								Status: JobStatusRunning,
							},
						},
						Total: 2,
					},
					StatusCode: http.StatusOK,
				},
			},
			want:    JobStatusFinished,
			wantErr: false,
		},
		{
			name:    "job running, then job done",
			storeID: "123",
			job:     DeployStoreJob,
			responses: []serverResponse{
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   DeployStoreJob,
								Status: JobStatusRunning,
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
								Type:   DeployStoreJob,
								Status: JobStatusFinished,
							},
						},
						Total: 1,
					},
					StatusCode: http.StatusOK,
				},
			},
			want:    JobStatusFinished,
			wantErr: false,
		},
		{
			name:    "job running, then job failed",
			storeID: "123",
			job:     DeployStoreJob,
			responses: []serverResponse{
				{
					Response: jobsResponse{
						Jobs: []jobsResponseJobItem{
							{
								JobID:  "456",
								Type:   DeployStoreJob,
								Status: JobStatusRunning,
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
								Type:   DeployStoreJob,
								Status: JobStatusErrored,
							},
						},
						Total: 1,
					},
					StatusCode: http.StatusOK,
				},
			},
			want:    JobStatusErrored,
			wantErr: false,
		},
		{
			name:    "http error",
			storeID: "123",
			job:     DeployStoreJob,
			responses: []serverResponse{
				{
					Response:   jobsResponse{},
					StatusCode: http.StatusGatewayTimeout,
				},
			},
			want:    JobStatusUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := 0

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/deployment/v2/data-stores/"+tt.storeID+"/jobs", r.URL.Path)

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

			svc := JobsClient{
				httpcli: NewTestHTTPClient(srv.URL),
				tick:    time.Second / 2,
				timeout: time.Second,
			}

			ctx := context.Background()

			got, err := svc.Await(ctx, tt.storeID, tt.job)
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
