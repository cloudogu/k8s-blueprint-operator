// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockSelfUpgradeUseCase is an autogenerated mock type for the selfUpgradeUseCase type
type mockSelfUpgradeUseCase struct {
	mock.Mock
}

type mockSelfUpgradeUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockSelfUpgradeUseCase) EXPECT() *mockSelfUpgradeUseCase_Expecter {
	return &mockSelfUpgradeUseCase_Expecter{mock: &_m.Mock}
}

// HandleSelfUpgrade provides a mock function with given fields: ctx, blueprintId
func (_m *mockSelfUpgradeUseCase) HandleSelfUpgrade(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for HandleSelfUpgrade")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSelfUpgradeUseCase_HandleSelfUpgrade_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HandleSelfUpgrade'
type mockSelfUpgradeUseCase_HandleSelfUpgrade_Call struct {
	*mock.Call
}

// HandleSelfUpgrade is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockSelfUpgradeUseCase_Expecter) HandleSelfUpgrade(ctx interface{}, blueprintId interface{}) *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call {
	return &mockSelfUpgradeUseCase_HandleSelfUpgrade_Call{Call: _e.mock.On("HandleSelfUpgrade", ctx, blueprintId)}
}

func (_c *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call) Run(run func(ctx context.Context, blueprintId string)) *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call) Return(_a0 error) *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call) RunAndReturn(run func(context.Context, string) error) *mockSelfUpgradeUseCase_HandleSelfUpgrade_Call {
	_c.Call.Return(run)
	return _c
}

// newMockSelfUpgradeUseCase creates a new instance of mockSelfUpgradeUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockSelfUpgradeUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockSelfUpgradeUseCase {
	mock := &mockSelfUpgradeUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
