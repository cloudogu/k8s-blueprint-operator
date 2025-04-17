// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	mock "github.com/stretchr/testify/mock"
)

// mockHealthConfigProvider is an autogenerated mock type for the healthConfigProvider type
type mockHealthConfigProvider struct {
	mock.Mock
}

type mockHealthConfigProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *mockHealthConfigProvider) EXPECT() *mockHealthConfigProvider_Expecter {
	return &mockHealthConfigProvider_Expecter{mock: &_m.Mock}
}

// GetRequiredComponents provides a mock function with given fields: ctx
func (_m *mockHealthConfigProvider) GetRequiredComponents(ctx context.Context) ([]ecosystem.RequiredComponent, error) {
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

// mockHealthConfigProvider_GetRequiredComponents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRequiredComponents'
type mockHealthConfigProvider_GetRequiredComponents_Call struct {
	*mock.Call
}

// GetRequiredComponents is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockHealthConfigProvider_Expecter) GetRequiredComponents(ctx interface{}) *mockHealthConfigProvider_GetRequiredComponents_Call {
	return &mockHealthConfigProvider_GetRequiredComponents_Call{Call: _e.mock.On("GetRequiredComponents", ctx)}
}

func (_c *mockHealthConfigProvider_GetRequiredComponents_Call) Run(run func(ctx context.Context)) *mockHealthConfigProvider_GetRequiredComponents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockHealthConfigProvider_GetRequiredComponents_Call) Return(_a0 []ecosystem.RequiredComponent, _a1 error) *mockHealthConfigProvider_GetRequiredComponents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockHealthConfigProvider_GetRequiredComponents_Call) RunAndReturn(run func(context.Context) ([]ecosystem.RequiredComponent, error)) *mockHealthConfigProvider_GetRequiredComponents_Call {
	_c.Call.Return(run)
	return _c
}

// GetWaitConfig provides a mock function with given fields: ctx
func (_m *mockHealthConfigProvider) GetWaitConfig(ctx context.Context) (ecosystem.WaitConfig, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetWaitConfig")
	}

	var r0 ecosystem.WaitConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (ecosystem.WaitConfig, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) ecosystem.WaitConfig); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(ecosystem.WaitConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockHealthConfigProvider_GetWaitConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWaitConfig'
type mockHealthConfigProvider_GetWaitConfig_Call struct {
	*mock.Call
}

// GetWaitConfig is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockHealthConfigProvider_Expecter) GetWaitConfig(ctx interface{}) *mockHealthConfigProvider_GetWaitConfig_Call {
	return &mockHealthConfigProvider_GetWaitConfig_Call{Call: _e.mock.On("GetWaitConfig", ctx)}
}

func (_c *mockHealthConfigProvider_GetWaitConfig_Call) Run(run func(ctx context.Context)) *mockHealthConfigProvider_GetWaitConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockHealthConfigProvider_GetWaitConfig_Call) Return(_a0 ecosystem.WaitConfig, _a1 error) *mockHealthConfigProvider_GetWaitConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockHealthConfigProvider_GetWaitConfig_Call) RunAndReturn(run func(context.Context) (ecosystem.WaitConfig, error)) *mockHealthConfigProvider_GetWaitConfig_Call {
	_c.Call.Return(run)
	return _c
}

// newMockHealthConfigProvider creates a new instance of mockHealthConfigProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockHealthConfigProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockHealthConfigProvider {
	mock := &mockHealthConfigProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
