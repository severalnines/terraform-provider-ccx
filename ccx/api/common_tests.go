package api

import (
	"context"
)

type fakeAuthorizer struct {
	wantToken string
	wantErr   error
}

func (f fakeAuthorizer) Auth(_ context.Context) (string, error) {
	return f.wantToken, f.wantErr
}
