// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockHttpClient is an autogenerated mock type for the HttpClient type
type MockHttpClient struct {
	mock.Mock
}

type MockHttpClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHttpClient) EXPECT() *MockHttpClient_Expecter {
	return &MockHttpClient_Expecter{mock: &_m.Mock}
}

// Do provides a mock function with given fields: ctx, method, path, body
func (_m *MockHttpClient) Do(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	ret := _m.Called(ctx, method, path, body)

	if len(ret) == 0 {
		panic("no return value specified for Do")
	}

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}) (*http.Response, error)); ok {
		return rf(ctx, method, path, body)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}) *http.Response); ok {
		r0 = rf(ctx, method, path, body)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, interface{}) error); ok {
		r1 = rf(ctx, method, path, body)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockHttpClient_Do_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Do'
type MockHttpClient_Do_Call struct {
	*mock.Call
}

// Do is a helper method to define mock.On call
//   - ctx context.Context
//   - method string
//   - path string
//   - body interface{}
func (_e *MockHttpClient_Expecter) Do(ctx interface{}, method interface{}, path interface{}, body interface{}) *MockHttpClient_Do_Call {
	return &MockHttpClient_Do_Call{Call: _e.mock.On("Do", ctx, method, path, body)}
}

func (_c *MockHttpClient_Do_Call) Run(run func(ctx context.Context, method string, path string, body interface{})) *MockHttpClient_Do_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(interface{}))
	})
	return _c
}

func (_c *MockHttpClient_Do_Call) Return(_a0 *http.Response, _a1 error) *MockHttpClient_Do_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHttpClient_Do_Call) RunAndReturn(run func(context.Context, string, string, interface{}) (*http.Response, error)) *MockHttpClient_Do_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, path, target
func (_m *MockHttpClient) Get(ctx context.Context, path string, target interface{}) error {
	ret := _m.Called(ctx, path, target)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, path, target)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockHttpClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockHttpClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - path string
//   - target interface{}
func (_e *MockHttpClient_Expecter) Get(ctx interface{}, path interface{}, target interface{}) *MockHttpClient_Get_Call {
	return &MockHttpClient_Get_Call{Call: _e.mock.On("Get", ctx, path, target)}
}

func (_c *MockHttpClient_Get_Call) Run(run func(ctx context.Context, path string, target interface{})) *MockHttpClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *MockHttpClient_Get_Call) Return(_a0 error) *MockHttpClient_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHttpClient_Get_Call) RunAndReturn(run func(context.Context, string, interface{}) error) *MockHttpClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHttpClient creates a new instance of MockHttpClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHttpClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHttpClient {
	mock := &MockHttpClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
