// Code generated by mockery v2.42.1. DO NOT EDIT.

package application

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockRegistryConfigUseCase is an autogenerated mock type for the registryConfigUseCase type
type mockRegistryConfigUseCase struct {
	mock.Mock
}

type mockRegistryConfigUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockRegistryConfigUseCase) EXPECT() *mockRegistryConfigUseCase_Expecter {
	return &mockRegistryConfigUseCase_Expecter{mock: &_m.Mock}
}

// ApplyConfig provides a mock function with given fields: ctx, blueprintId
func (_m *mockRegistryConfigUseCase) ApplyConfig(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for ApplyConfig")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockRegistryConfigUseCase_ApplyConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyConfig'
type mockRegistryConfigUseCase_ApplyConfig_Call struct {
	*mock.Call
}

// ApplyConfig is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockRegistryConfigUseCase_Expecter) ApplyConfig(ctx interface{}, blueprintId interface{}) *mockRegistryConfigUseCase_ApplyConfig_Call {
	return &mockRegistryConfigUseCase_ApplyConfig_Call{Call: _e.mock.On("ApplyConfig", ctx, blueprintId)}
}

func (_c *mockRegistryConfigUseCase_ApplyConfig_Call) Run(run func(ctx context.Context, blueprintId string)) *mockRegistryConfigUseCase_ApplyConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockRegistryConfigUseCase_ApplyConfig_Call) Return(_a0 error) *mockRegistryConfigUseCase_ApplyConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistryConfigUseCase_ApplyConfig_Call) RunAndReturn(run func(context.Context, string) error) *mockRegistryConfigUseCase_ApplyConfig_Call {
	_c.Call.Return(run)
	return _c
}

// newMockRegistryConfigUseCase creates a new instance of mockRegistryConfigUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRegistryConfigUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRegistryConfigUseCase {
	mock := &mockRegistryConfigUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
