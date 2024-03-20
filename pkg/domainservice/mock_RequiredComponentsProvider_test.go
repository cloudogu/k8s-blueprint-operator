// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	mock "github.com/stretchr/testify/mock"
)

// MockRequiredComponentsProvider is an autogenerated mock type for the RequiredComponentsProvider type
type MockRequiredComponentsProvider struct {
	mock.Mock
}

type MockRequiredComponentsProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRequiredComponentsProvider) EXPECT() *MockRequiredComponentsProvider_Expecter {
	return &MockRequiredComponentsProvider_Expecter{mock: &_m.Mock}
}

// GetRequiredComponents provides a mock function with given fields: ctx
func (_m *MockRequiredComponentsProvider) GetRequiredComponents(ctx context.Context) ([]ecosystem.RequiredComponent, error) {
	ret := _m.Called(ctx)

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

// MockRequiredComponentsProvider_GetRequiredComponents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRequiredComponents'
type MockRequiredComponentsProvider_GetRequiredComponents_Call struct {
	*mock.Call
}

// GetRequiredComponents is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockRequiredComponentsProvider_Expecter) GetRequiredComponents(ctx interface{}) *MockRequiredComponentsProvider_GetRequiredComponents_Call {
	return &MockRequiredComponentsProvider_GetRequiredComponents_Call{Call: _e.mock.On("GetRequiredComponents", ctx)}
}

func (_c *MockRequiredComponentsProvider_GetRequiredComponents_Call) Run(run func(ctx context.Context)) *MockRequiredComponentsProvider_GetRequiredComponents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockRequiredComponentsProvider_GetRequiredComponents_Call) Return(_a0 []ecosystem.RequiredComponent, _a1 error) *MockRequiredComponentsProvider_GetRequiredComponents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRequiredComponentsProvider_GetRequiredComponents_Call) RunAndReturn(run func(context.Context) ([]ecosystem.RequiredComponent, error)) *MockRequiredComponentsProvider_GetRequiredComponents_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockRequiredComponentsProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockRequiredComponentsProvider creates a new instance of MockRequiredComponentsProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockRequiredComponentsProvider(t mockConstructorTestingTNewMockRequiredComponentsProvider) *MockRequiredComponentsProvider {
	mock := &MockRequiredComponentsProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
