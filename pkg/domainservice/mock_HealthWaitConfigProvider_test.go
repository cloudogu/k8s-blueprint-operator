// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	mock "github.com/stretchr/testify/mock"
)

// MockHealthWaitConfigProvider is an autogenerated mock type for the HealthWaitConfigProvider type
type MockHealthWaitConfigProvider struct {
	mock.Mock
}

type MockHealthWaitConfigProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHealthWaitConfigProvider) EXPECT() *MockHealthWaitConfigProvider_Expecter {
	return &MockHealthWaitConfigProvider_Expecter{mock: &_m.Mock}
}

// GetWaitConfig provides a mock function with given fields: ctx
func (_m *MockHealthWaitConfigProvider) GetWaitConfig(ctx context.Context) (ecosystem.WaitConfig, error) {
	ret := _m.Called(ctx)

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

// MockHealthWaitConfigProvider_GetWaitConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetWaitConfig'
type MockHealthWaitConfigProvider_GetWaitConfig_Call struct {
	*mock.Call
}

// GetWaitConfig is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockHealthWaitConfigProvider_Expecter) GetWaitConfig(ctx interface{}) *MockHealthWaitConfigProvider_GetWaitConfig_Call {
	return &MockHealthWaitConfigProvider_GetWaitConfig_Call{Call: _e.mock.On("GetWaitConfig", ctx)}
}

func (_c *MockHealthWaitConfigProvider_GetWaitConfig_Call) Run(run func(ctx context.Context)) *MockHealthWaitConfigProvider_GetWaitConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockHealthWaitConfigProvider_GetWaitConfig_Call) Return(_a0 ecosystem.WaitConfig, _a1 error) *MockHealthWaitConfigProvider_GetWaitConfig_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockHealthWaitConfigProvider_GetWaitConfig_Call) RunAndReturn(run func(context.Context) (ecosystem.WaitConfig, error)) *MockHealthWaitConfigProvider_GetWaitConfig_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockHealthWaitConfigProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockHealthWaitConfigProvider creates a new instance of MockHealthWaitConfigProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockHealthWaitConfigProvider(t mockConstructorTestingTNewMockHealthWaitConfigProvider) *MockHealthWaitConfigProvider {
	mock := &MockHealthWaitConfigProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
