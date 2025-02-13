// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	ccx "github.com/severalnines/terraform-provider-ccx/internal/ccx"

	mock "github.com/stretchr/testify/mock"
)

// MockContentService is an autogenerated mock type for the ContentService type
type MockContentService struct {
	mock.Mock
}

type MockContentService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockContentService) EXPECT() *MockContentService_Expecter {
	return &MockContentService_Expecter{mock: &_m.Mock}
}

// AvailabilityZones provides a mock function with given fields: ctx, provider, region
func (_m *MockContentService) AvailabilityZones(ctx context.Context, provider string, region string) ([]string, error) {
	ret := _m.Called(ctx, provider, region)

	if len(ret) == 0 {
		panic("no return value specified for AvailabilityZones")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]string, error)); ok {
		return rf(ctx, provider, region)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []string); ok {
		r0 = rf(ctx, provider, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, provider, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContentService_AvailabilityZones_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AvailabilityZones'
type MockContentService_AvailabilityZones_Call struct {
	*mock.Call
}

// AvailabilityZones is a helper method to define mock.On call
//   - ctx context.Context
//   - provider string
//   - region string
func (_e *MockContentService_Expecter) AvailabilityZones(ctx interface{}, provider interface{}, region interface{}) *MockContentService_AvailabilityZones_Call {
	return &MockContentService_AvailabilityZones_Call{Call: _e.mock.On("AvailabilityZones", ctx, provider, region)}
}

func (_c *MockContentService_AvailabilityZones_Call) Run(run func(ctx context.Context, provider string, region string)) *MockContentService_AvailabilityZones_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockContentService_AvailabilityZones_Call) Return(_a0 []string, _a1 error) *MockContentService_AvailabilityZones_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContentService_AvailabilityZones_Call) RunAndReturn(run func(context.Context, string, string) ([]string, error)) *MockContentService_AvailabilityZones_Call {
	_c.Call.Return(run)
	return _c
}

// DBVendors provides a mock function with given fields: ctx
func (_m *MockContentService) DBVendors(ctx context.Context) ([]ccx.DBVendorInfo, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DBVendors")
	}

	var r0 []ccx.DBVendorInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]ccx.DBVendorInfo, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []ccx.DBVendorInfo); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ccx.DBVendorInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContentService_DBVendors_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DBVendors'
type MockContentService_DBVendors_Call struct {
	*mock.Call
}

// DBVendors is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockContentService_Expecter) DBVendors(ctx interface{}) *MockContentService_DBVendors_Call {
	return &MockContentService_DBVendors_Call{Call: _e.mock.On("DBVendors", ctx)}
}

func (_c *MockContentService_DBVendors_Call) Run(run func(ctx context.Context)) *MockContentService_DBVendors_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockContentService_DBVendors_Call) Return(_a0 []ccx.DBVendorInfo, _a1 error) *MockContentService_DBVendors_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContentService_DBVendors_Call) RunAndReturn(run func(context.Context) ([]ccx.DBVendorInfo, error)) *MockContentService_DBVendors_Call {
	_c.Call.Return(run)
	return _c
}

// InstanceSizes provides a mock function with given fields: ctx
func (_m *MockContentService) InstanceSizes(ctx context.Context) (map[string][]ccx.InstanceSize, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for InstanceSizes")
	}

	var r0 map[string][]ccx.InstanceSize
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[string][]ccx.InstanceSize, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[string][]ccx.InstanceSize); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]ccx.InstanceSize)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContentService_InstanceSizes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InstanceSizes'
type MockContentService_InstanceSizes_Call struct {
	*mock.Call
}

// InstanceSizes is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockContentService_Expecter) InstanceSizes(ctx interface{}) *MockContentService_InstanceSizes_Call {
	return &MockContentService_InstanceSizes_Call{Call: _e.mock.On("InstanceSizes", ctx)}
}

func (_c *MockContentService_InstanceSizes_Call) Run(run func(ctx context.Context)) *MockContentService_InstanceSizes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockContentService_InstanceSizes_Call) Return(_a0 map[string][]ccx.InstanceSize, _a1 error) *MockContentService_InstanceSizes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContentService_InstanceSizes_Call) RunAndReturn(run func(context.Context) (map[string][]ccx.InstanceSize, error)) *MockContentService_InstanceSizes_Call {
	_c.Call.Return(run)
	return _c
}

// VolumeTypes provides a mock function with given fields: ctx, cloud
func (_m *MockContentService) VolumeTypes(ctx context.Context, cloud string) ([]string, error) {
	ret := _m.Called(ctx, cloud)

	if len(ret) == 0 {
		panic("no return value specified for VolumeTypes")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(ctx, cloud)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, cloud)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, cloud)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockContentService_VolumeTypes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VolumeTypes'
type MockContentService_VolumeTypes_Call struct {
	*mock.Call
}

// VolumeTypes is a helper method to define mock.On call
//   - ctx context.Context
//   - cloud string
func (_e *MockContentService_Expecter) VolumeTypes(ctx interface{}, cloud interface{}) *MockContentService_VolumeTypes_Call {
	return &MockContentService_VolumeTypes_Call{Call: _e.mock.On("VolumeTypes", ctx, cloud)}
}

func (_c *MockContentService_VolumeTypes_Call) Run(run func(ctx context.Context, cloud string)) *MockContentService_VolumeTypes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockContentService_VolumeTypes_Call) Return(_a0 []string, _a1 error) *MockContentService_VolumeTypes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockContentService_VolumeTypes_Call) RunAndReturn(run func(context.Context, string) ([]string, error)) *MockContentService_VolumeTypes_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockContentService creates a new instance of MockContentService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockContentService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockContentService {
	mock := &MockContentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
