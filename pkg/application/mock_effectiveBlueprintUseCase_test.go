// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockEffectiveBlueprintUseCase is an autogenerated mock type for the effectiveBlueprintUseCase type
type mockEffectiveBlueprintUseCase struct {
	mock.Mock
}

type mockEffectiveBlueprintUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockEffectiveBlueprintUseCase) EXPECT() *mockEffectiveBlueprintUseCase_Expecter {
	return &mockEffectiveBlueprintUseCase_Expecter{mock: &_m.Mock}
}

// CalculateEffectiveBlueprint provides a mock function with given fields: ctx, blueprintId
func (_m *mockEffectiveBlueprintUseCase) CalculateEffectiveBlueprint(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for CalculateEffectiveBlueprint")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CalculateEffectiveBlueprint'
type mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call struct {
	*mock.Call
}

// CalculateEffectiveBlueprint is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockEffectiveBlueprintUseCase_Expecter) CalculateEffectiveBlueprint(ctx interface{}, blueprintId interface{}) *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call {
	return &mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call{Call: _e.mock.On("CalculateEffectiveBlueprint", ctx, blueprintId)}
}

func (_c *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call) Run(run func(ctx context.Context, blueprintId string)) *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call) Return(_a0 error) *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call) RunAndReturn(run func(context.Context, string) error) *mockEffectiveBlueprintUseCase_CalculateEffectiveBlueprint_Call {
	_c.Call.Return(run)
	return _c
}

// newMockEffectiveBlueprintUseCase creates a new instance of mockEffectiveBlueprintUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockEffectiveBlueprintUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockEffectiveBlueprintUseCase {
	mock := &mockEffectiveBlueprintUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
