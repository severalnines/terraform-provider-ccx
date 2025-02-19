// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	ccx "github.com/severalnines/terraform-provider-ccx/internal/ccx"

	mock "github.com/stretchr/testify/mock"
)

// MockDatastoreService is an autogenerated mock type for the DatastoreService type
type MockDatastoreService struct {
	mock.Mock
}

type MockDatastoreService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatastoreService) EXPECT() *MockDatastoreService_Expecter {
	return &MockDatastoreService_Expecter{mock: &_m.Mock}
}

// ApplyParameterGroup provides a mock function with given fields: ctx, id, group
func (_m *MockDatastoreService) ApplyParameterGroup(ctx context.Context, id string, group string) error {
	ret := _m.Called(ctx, id, group)

	if len(ret) == 0 {
		panic("no return value specified for ApplyParameterGroup")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, id, group)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatastoreService_ApplyParameterGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyParameterGroup'
type MockDatastoreService_ApplyParameterGroup_Call struct {
	*mock.Call
}

// ApplyParameterGroup is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
//   - group string
func (_e *MockDatastoreService_Expecter) ApplyParameterGroup(ctx interface{}, id interface{}, group interface{}) *MockDatastoreService_ApplyParameterGroup_Call {
	return &MockDatastoreService_ApplyParameterGroup_Call{Call: _e.mock.On("ApplyParameterGroup", ctx, id, group)}
}

func (_c *MockDatastoreService_ApplyParameterGroup_Call) Run(run func(ctx context.Context, id string, group string)) *MockDatastoreService_ApplyParameterGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDatastoreService_ApplyParameterGroup_Call) Return(_a0 error) *MockDatastoreService_ApplyParameterGroup_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatastoreService_ApplyParameterGroup_Call) RunAndReturn(run func(context.Context, string, string) error) *MockDatastoreService_ApplyParameterGroup_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function with given fields: ctx, c
func (_m *MockDatastoreService) Create(ctx context.Context, c ccx.Datastore) (*ccx.Datastore, error) {
	ret := _m.Called(ctx, c)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *ccx.Datastore
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ccx.Datastore) (*ccx.Datastore, error)); ok {
		return rf(ctx, c)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ccx.Datastore) *ccx.Datastore); ok {
		r0 = rf(ctx, c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ccx.Datastore)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ccx.Datastore) error); ok {
		r1 = rf(ctx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatastoreService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockDatastoreService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - c ccx.Datastore
func (_e *MockDatastoreService_Expecter) Create(ctx interface{}, c interface{}) *MockDatastoreService_Create_Call {
	return &MockDatastoreService_Create_Call{Call: _e.mock.On("Create", ctx, c)}
}

func (_c *MockDatastoreService_Create_Call) Run(run func(ctx context.Context, c ccx.Datastore)) *MockDatastoreService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ccx.Datastore))
	})
	return _c
}

func (_c *MockDatastoreService_Create_Call) Return(_a0 *ccx.Datastore, _a1 error) *MockDatastoreService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatastoreService_Create_Call) RunAndReturn(run func(context.Context, ccx.Datastore) (*ccx.Datastore, error)) *MockDatastoreService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MockDatastoreService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatastoreService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockDatastoreService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockDatastoreService_Expecter) Delete(ctx interface{}, id interface{}) *MockDatastoreService_Delete_Call {
	return &MockDatastoreService_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockDatastoreService_Delete_Call) Run(run func(ctx context.Context, id string)) *MockDatastoreService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatastoreService_Delete_Call) Return(_a0 error) *MockDatastoreService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatastoreService_Delete_Call) RunAndReturn(run func(context.Context, string) error) *MockDatastoreService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: ctx, id
func (_m *MockDatastoreService) Read(ctx context.Context, id string) (*ccx.Datastore, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 *ccx.Datastore
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*ccx.Datastore, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *ccx.Datastore); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ccx.Datastore)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatastoreService_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockDatastoreService_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - ctx context.Context
//   - id string
func (_e *MockDatastoreService_Expecter) Read(ctx interface{}, id interface{}) *MockDatastoreService_Read_Call {
	return &MockDatastoreService_Read_Call{Call: _e.mock.On("Read", ctx, id)}
}

func (_c *MockDatastoreService_Read_Call) Run(run func(ctx context.Context, id string)) *MockDatastoreService_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatastoreService_Read_Call) Return(_a0 *ccx.Datastore, _a1 error) *MockDatastoreService_Read_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatastoreService_Read_Call) RunAndReturn(run func(context.Context, string) (*ccx.Datastore, error)) *MockDatastoreService_Read_Call {
	_c.Call.Return(run)
	return _c
}

// SetFirewallRules provides a mock function with given fields: ctx, storeID, firewalls
func (_m *MockDatastoreService) SetFirewallRules(ctx context.Context, storeID string, firewalls []ccx.FirewallRule) error {
	ret := _m.Called(ctx, storeID, firewalls)

	if len(ret) == 0 {
		panic("no return value specified for SetFirewallRules")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []ccx.FirewallRule) error); ok {
		r0 = rf(ctx, storeID, firewalls)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatastoreService_SetFirewallRules_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetFirewallRules'
type MockDatastoreService_SetFirewallRules_Call struct {
	*mock.Call
}

// SetFirewallRules is a helper method to define mock.On call
//   - ctx context.Context
//   - storeID string
//   - firewalls []ccx.FirewallRule
func (_e *MockDatastoreService_Expecter) SetFirewallRules(ctx interface{}, storeID interface{}, firewalls interface{}) *MockDatastoreService_SetFirewallRules_Call {
	return &MockDatastoreService_SetFirewallRules_Call{Call: _e.mock.On("SetFirewallRules", ctx, storeID, firewalls)}
}

func (_c *MockDatastoreService_SetFirewallRules_Call) Run(run func(ctx context.Context, storeID string, firewalls []ccx.FirewallRule)) *MockDatastoreService_SetFirewallRules_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]ccx.FirewallRule))
	})
	return _c
}

func (_c *MockDatastoreService_SetFirewallRules_Call) Return(_a0 error) *MockDatastoreService_SetFirewallRules_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatastoreService_SetFirewallRules_Call) RunAndReturn(run func(context.Context, string, []ccx.FirewallRule) error) *MockDatastoreService_SetFirewallRules_Call {
	_c.Call.Return(run)
	return _c
}

// SetMaintenanceSettings provides a mock function with given fields: ctx, storeID, settings
func (_m *MockDatastoreService) SetMaintenanceSettings(ctx context.Context, storeID string, settings ccx.MaintenanceSettings) error {
	ret := _m.Called(ctx, storeID, settings)

	if len(ret) == 0 {
		panic("no return value specified for SetMaintenanceSettings")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ccx.MaintenanceSettings) error); ok {
		r0 = rf(ctx, storeID, settings)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatastoreService_SetMaintenanceSettings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetMaintenanceSettings'
type MockDatastoreService_SetMaintenanceSettings_Call struct {
	*mock.Call
}

// SetMaintenanceSettings is a helper method to define mock.On call
//   - ctx context.Context
//   - storeID string
//   - settings ccx.MaintenanceSettings
func (_e *MockDatastoreService_Expecter) SetMaintenanceSettings(ctx interface{}, storeID interface{}, settings interface{}) *MockDatastoreService_SetMaintenanceSettings_Call {
	return &MockDatastoreService_SetMaintenanceSettings_Call{Call: _e.mock.On("SetMaintenanceSettings", ctx, storeID, settings)}
}

func (_c *MockDatastoreService_SetMaintenanceSettings_Call) Run(run func(ctx context.Context, storeID string, settings ccx.MaintenanceSettings)) *MockDatastoreService_SetMaintenanceSettings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(ccx.MaintenanceSettings))
	})
	return _c
}

func (_c *MockDatastoreService_SetMaintenanceSettings_Call) Return(_a0 error) *MockDatastoreService_SetMaintenanceSettings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatastoreService_SetMaintenanceSettings_Call) RunAndReturn(run func(context.Context, string, ccx.MaintenanceSettings) error) *MockDatastoreService_SetMaintenanceSettings_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, old, next
func (_m *MockDatastoreService) Update(ctx context.Context, old ccx.Datastore, next ccx.Datastore) (*ccx.Datastore, error) {
	ret := _m.Called(ctx, old, next)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *ccx.Datastore
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ccx.Datastore, ccx.Datastore) (*ccx.Datastore, error)); ok {
		return rf(ctx, old, next)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ccx.Datastore, ccx.Datastore) *ccx.Datastore); ok {
		r0 = rf(ctx, old, next)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ccx.Datastore)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ccx.Datastore, ccx.Datastore) error); ok {
		r1 = rf(ctx, old, next)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatastoreService_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockDatastoreService_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - old ccx.Datastore
//   - next ccx.Datastore
func (_e *MockDatastoreService_Expecter) Update(ctx interface{}, old interface{}, next interface{}) *MockDatastoreService_Update_Call {
	return &MockDatastoreService_Update_Call{Call: _e.mock.On("Update", ctx, old, next)}
}

func (_c *MockDatastoreService_Update_Call) Run(run func(ctx context.Context, old ccx.Datastore, next ccx.Datastore)) *MockDatastoreService_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ccx.Datastore), args[2].(ccx.Datastore))
	})
	return _c
}

func (_c *MockDatastoreService_Update_Call) Return(_a0 *ccx.Datastore, _a1 error) *MockDatastoreService_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatastoreService_Update_Call) RunAndReturn(run func(context.Context, ccx.Datastore, ccx.Datastore) (*ccx.Datastore, error)) *MockDatastoreService_Update_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatastoreService creates a new instance of MockDatastoreService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatastoreService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatastoreService {
	mock := &MockDatastoreService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
