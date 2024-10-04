// Code generated by mockery v2.42.1. DO NOT EDIT.

package application

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockDoguRestartUseCase is an autogenerated mock type for the doguRestartUseCase type
type mockDoguRestartUseCase struct {
	mock.Mock
}

type mockDoguRestartUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguRestartUseCase) EXPECT() *mockDoguRestartUseCase_Expecter {
	return &mockDoguRestartUseCase_Expecter{mock: &_m.Mock}
}

// TriggerDoguRestarts provides a mock function with given fields: ctx, blueprintid
func (_m *mockDoguRestartUseCase) TriggerDoguRestarts(ctx context.Context, blueprintid string) error {
	ret := _m.Called(ctx, blueprintid)

	if len(ret) == 0 {
		panic("no return value specified for TriggerDoguRestarts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguRestartUseCase_TriggerDoguRestarts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TriggerDoguRestarts'
type mockDoguRestartUseCase_TriggerDoguRestarts_Call struct {
	*mock.Call
}

// TriggerDoguRestarts is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintid string
func (_e *mockDoguRestartUseCase_Expecter) TriggerDoguRestarts(ctx interface{}, blueprintid interface{}) *mockDoguRestartUseCase_TriggerDoguRestarts_Call {
	return &mockDoguRestartUseCase_TriggerDoguRestarts_Call{Call: _e.mock.On("TriggerDoguRestarts", ctx, blueprintid)}
}

func (_c *mockDoguRestartUseCase_TriggerDoguRestarts_Call) Run(run func(ctx context.Context, blueprintid string)) *mockDoguRestartUseCase_TriggerDoguRestarts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockDoguRestartUseCase_TriggerDoguRestarts_Call) Return(_a0 error) *mockDoguRestartUseCase_TriggerDoguRestarts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguRestartUseCase_TriggerDoguRestarts_Call) RunAndReturn(run func(context.Context, string) error) *mockDoguRestartUseCase_TriggerDoguRestarts_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguRestartUseCase creates a new instance of mockDoguRestartUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguRestartUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguRestartUseCase {
	mock := &mockDoguRestartUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
