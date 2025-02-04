package mocks

import (
	"github.com/stretchr/testify/mock"
)

func MockHttpClientExpectGet[T any](m *MockHttpClient, path string, r func(t *T) bool, wantErr error) {
	m.EXPECT().Get(mock.Anything, path, mock.MatchedBy(func(rs *T) bool {
		return r(rs)
	})).Return(wantErr)
}
