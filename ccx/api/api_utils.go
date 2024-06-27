package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/severalnines/terraform-provider-ccx/ccx"
	"github.com/severalnines/terraform-provider-ccx/internal/lib"
)

// pString returns the value of a string pointer or empty string if pointer is nil
func pString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

// pUint64 returns the value of a uint64 pointer or empty uint64 if pointer is nil
func pUint64(n *uint64) uint64 {
	if n == nil {
		return 0
	}

	return *n
}

// httpGet is a helper function to send a GET request to the API and decode the json response into target
func httpGet(ctx context.Context, auth authorizer, url string, target any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Join(ccx.RequestInitializationErr, err)
	}

	token, err := auth.Auth(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	client := &http.Client{Timeout: ccx.DefaultTimeout}

	res, err := client.Do(req)
	if err != nil {
		return errors.Join(ccx.RequestSendingErr, err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status = %d", lib.ErrorFromErrorResponse(res.Body), res.StatusCode)
	}

	defer func() {
		if res.Body != nil {
			_ = res.Body.Close()
		}
	}()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.Join(ccx.ResponseReadFailedErr, err)
	}

	if err := json.Unmarshal(b, target); err != nil {
		return errors.Join(ccx.ResponseDecodingErr, err)
	}

	return nil
}
