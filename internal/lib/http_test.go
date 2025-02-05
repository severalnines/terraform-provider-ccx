package lib

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name    string
		Code    json.Number
		Err     string
		ErrLong string
		status  int
		want    string
	}{
		{
			name:    "Error with code and short error",
			Code:    "1337",
			Err:     "datastore does not exist",
			ErrLong: "",
			status:  http.StatusNotFound,
			want:    "datastore does not exist (code: 1337, response: 404 - Not Found)",
		},
		{
			name:    "Error with long error",
			Code:    "",
			Err:     "",
			ErrLong: "failed to connect to database",
			status:  http.StatusInternalServerError,
			want:    "failed to connect to database (response: 500 - Internal Server Error)",
		},
		{
			name:    "Error with code and long error",
			Code:    "B9d",
			Err:     "",
			ErrLong: "something went wrong",
			status:  http.StatusInternalServerError,
			want:    "something went wrong (code: B9d, response: 500 - Internal Server Error)",
		},
		{
			name:    "Error with no code and no error message",
			Code:    "",
			Err:     "",
			ErrLong: "",
			status:  http.StatusUnauthorized,
			want:    "an error occurred (response: 401 - Unauthorized)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ErrorResponse{
				Code:    tt.Code,
				Err:     tt.Err,
				ErrLong: tt.ErrLong,
				status:  tt.status,
			}
			assert.Equalf(t, tt.want, r.Error(), "Error()")
		})
	}
}
