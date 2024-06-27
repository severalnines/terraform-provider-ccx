package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/severalnines/terraform-provider-ccx/ccx"
)

// ErrorResponse represents generic error responses from ccx api
type ErrorResponse struct {
	Code    json.Number `json:"code"`
	Err     string      `json:"err"`
	ErrLong string      `json:"error"`
}

func (r ErrorResponse) Error() string {
	prefix := r.Code.String()
	if prefix != "" {
		prefix += ": "
	}
	if r.Err != "" {
		return prefix + r.Err
	}
	return prefix + r.ErrLong
}

// ErrorFromErrorResponse decodes the body of type ErrorResponse and returns an error
func ErrorFromErrorResponse(body io.ReadCloser) error {
	defer Closed(body)
	b, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("%w: reason could not be read", err)
	}

	var r ErrorResponse
	err = json.Unmarshal(b, &r)
	if err != nil {
		return fmt.Errorf("%w: reason could not be decoded", err)
	}

	return r
}

// DecodeJsonInto is a helper to decode JSON body into a target type
func DecodeJsonInto(body io.ReadCloser, target any) error {
	defer Closed(body)

	raw, err := io.ReadAll(body)
	if err != nil {
		return errors.Join(ccx.ResponseReadFailedErr, err)
	}

	err = json.Unmarshal(raw, target)
	if err != nil {
		return errors.Join(ccx.ResponseDecodingErr, err)
	}

	return nil
}

// Closed will close an io.Closer object silently
// useful for deferred closing
func Closed(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
