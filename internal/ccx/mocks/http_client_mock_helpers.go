package mocks

import (
	"github.com/stretchr/testify/mock"
)

func MockHttpClientExpectGet[T any](m *MockHttpClient, path string, t T, wantErr error) {
	m.EXPECT().Get(mock.Anything, path, mock.MatchedBy(func(rs *T) bool {
		*rs = t
		return true
	})).Return(wantErr)
}
