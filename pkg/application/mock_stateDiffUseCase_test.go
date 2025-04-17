// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockStateDiffUseCase is an autogenerated mock type for the stateDiffUseCase type
type mockStateDiffUseCase struct {
	mock.Mock
}

type mockStateDiffUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockStateDiffUseCase) EXPECT() *mockStateDiffUseCase_Expecter {
	return &mockStateDiffUseCase_Expecter{mock: &_m.Mock}
}

// DetermineStateDiff provides a mock function with given fields: ctx, blueprintId
func (_m *mockStateDiffUseCase) DetermineStateDiff(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for DetermineStateDiff")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockStateDiffUseCase_DetermineStateDiff_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DetermineStateDiff'
type mockStateDiffUseCase_DetermineStateDiff_Call struct {
	*mock.Call
}

// DetermineStateDiff is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockStateDiffUseCase_Expecter) DetermineStateDiff(ctx interface{}, blueprintId interface{}) *mockStateDiffUseCase_DetermineStateDiff_Call {
	return &mockStateDiffUseCase_DetermineStateDiff_Call{Call: _e.mock.On("DetermineStateDiff", ctx, blueprintId)}
}

func (_c *mockStateDiffUseCase_DetermineStateDiff_Call) Run(run func(ctx context.Context, blueprintId string)) *mockStateDiffUseCase_DetermineStateDiff_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockStateDiffUseCase_DetermineStateDiff_Call) Return(_a0 error) *mockStateDiffUseCase_DetermineStateDiff_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockStateDiffUseCase_DetermineStateDiff_Call) RunAndReturn(run func(context.Context, string) error) *mockStateDiffUseCase_DetermineStateDiff_Call {
	_c.Call.Return(run)
	return _c
}

// newMockStateDiffUseCase creates a new instance of mockStateDiffUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockStateDiffUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockStateDiffUseCase {
	mock := &mockStateDiffUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
