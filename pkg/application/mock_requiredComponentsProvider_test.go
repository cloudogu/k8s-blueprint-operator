// Code generated by mockery v2.42.1. DO NOT EDIT.

package application

import (
	context "context"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	mock "github.com/stretchr/testify/mock"
)

// mockRequiredComponentsProvider is an autogenerated mock type for the requiredComponentsProvider type
type mockRequiredComponentsProvider struct {
	mock.Mock
}

type mockRequiredComponentsProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *mockRequiredComponentsProvider) EXPECT() *mockRequiredComponentsProvider_Expecter {
	return &mockRequiredComponentsProvider_Expecter{mock: &_m.Mock}
}

// GetRequiredComponents provides a mock function with given fields: ctx
func (_m *mockRequiredComponentsProvider) GetRequiredComponents(ctx context.Context) ([]ecosystem.RequiredComponent, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetRequiredComponents")
	}

	var r0 []ecosystem.RequiredComponent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]ecosystem.RequiredComponent, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []ecosystem.RequiredComponent); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ecosystem.RequiredComponent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockRequiredComponentsProvider_GetRequiredComponents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRequiredComponents'
type mockRequiredComponentsProvider_GetRequiredComponents_Call struct {
	*mock.Call
}

// GetRequiredComponents is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockRequiredComponentsProvider_Expecter) GetRequiredComponents(ctx interface{}) *mockRequiredComponentsProvider_GetRequiredComponents_Call {
	return &mockRequiredComponentsProvider_GetRequiredComponents_Call{Call: _e.mock.On("GetRequiredComponents", ctx)}
}

func (_c *mockRequiredComponentsProvider_GetRequiredComponents_Call) Run(run func(ctx context.Context)) *mockRequiredComponentsProvider_GetRequiredComponents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockRequiredComponentsProvider_GetRequiredComponents_Call) Return(_a0 []ecosystem.RequiredComponent, _a1 error) *mockRequiredComponentsProvider_GetRequiredComponents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockRequiredComponentsProvider_GetRequiredComponents_Call) RunAndReturn(run func(context.Context) ([]ecosystem.RequiredComponent, error)) *mockRequiredComponentsProvider_GetRequiredComponents_Call {
	_c.Call.Return(run)
	return _c
}

// newMockRequiredComponentsProvider creates a new instance of mockRequiredComponentsProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRequiredComponentsProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRequiredComponentsProvider {
	mock := &mockRequiredComponentsProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
