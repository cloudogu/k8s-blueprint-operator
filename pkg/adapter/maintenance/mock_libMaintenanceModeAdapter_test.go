// Code generated by mockery v2.53.3. DO NOT EDIT.

package maintenance

import (
	context "context"

	repository "github.com/cloudogu/k8s-registry-lib/repository"
	mock "github.com/stretchr/testify/mock"
)

// mockLibMaintenanceModeAdapter is an autogenerated mock type for the libMaintenanceModeAdapter type
type mockLibMaintenanceModeAdapter struct {
	mock.Mock
}

type mockLibMaintenanceModeAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockLibMaintenanceModeAdapter) EXPECT() *mockLibMaintenanceModeAdapter_Expecter {
	return &mockLibMaintenanceModeAdapter_Expecter{mock: &_m.Mock}
}

// Activate provides a mock function with given fields: ctx, content
func (_m *mockLibMaintenanceModeAdapter) Activate(ctx context.Context, content repository.MaintenanceModeDescription) error {
	ret := _m.Called(ctx, content)

	if len(ret) == 0 {
		panic("no return value specified for Activate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.MaintenanceModeDescription) error); ok {
		r0 = rf(ctx, content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockLibMaintenanceModeAdapter_Activate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Activate'
type mockLibMaintenanceModeAdapter_Activate_Call struct {
	*mock.Call
}

// Activate is a helper method to define mock.On call
//   - ctx context.Context
//   - content repository.MaintenanceModeDescription
func (_e *mockLibMaintenanceModeAdapter_Expecter) Activate(ctx interface{}, content interface{}) *mockLibMaintenanceModeAdapter_Activate_Call {
	return &mockLibMaintenanceModeAdapter_Activate_Call{Call: _e.mock.On("Activate", ctx, content)}
}

func (_c *mockLibMaintenanceModeAdapter_Activate_Call) Run(run func(ctx context.Context, content repository.MaintenanceModeDescription)) *mockLibMaintenanceModeAdapter_Activate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repository.MaintenanceModeDescription))
	})
	return _c
}

func (_c *mockLibMaintenanceModeAdapter_Activate_Call) Return(_a0 error) *mockLibMaintenanceModeAdapter_Activate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockLibMaintenanceModeAdapter_Activate_Call) RunAndReturn(run func(context.Context, repository.MaintenanceModeDescription) error) *mockLibMaintenanceModeAdapter_Activate_Call {
	_c.Call.Return(run)
	return _c
}

// Deactivate provides a mock function with given fields: ctx
func (_m *mockLibMaintenanceModeAdapter) Deactivate(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Deactivate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockLibMaintenanceModeAdapter_Deactivate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deactivate'
type mockLibMaintenanceModeAdapter_Deactivate_Call struct {
	*mock.Call
}

// Deactivate is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockLibMaintenanceModeAdapter_Expecter) Deactivate(ctx interface{}) *mockLibMaintenanceModeAdapter_Deactivate_Call {
	return &mockLibMaintenanceModeAdapter_Deactivate_Call{Call: _e.mock.On("Deactivate", ctx)}
}

func (_c *mockLibMaintenanceModeAdapter_Deactivate_Call) Run(run func(ctx context.Context)) *mockLibMaintenanceModeAdapter_Deactivate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockLibMaintenanceModeAdapter_Deactivate_Call) Return(_a0 error) *mockLibMaintenanceModeAdapter_Deactivate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockLibMaintenanceModeAdapter_Deactivate_Call) RunAndReturn(run func(context.Context) error) *mockLibMaintenanceModeAdapter_Deactivate_Call {
	_c.Call.Return(run)
	return _c
}

// newMockLibMaintenanceModeAdapter creates a new instance of mockLibMaintenanceModeAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockLibMaintenanceModeAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockLibMaintenanceModeAdapter {
	mock := &mockLibMaintenanceModeAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
