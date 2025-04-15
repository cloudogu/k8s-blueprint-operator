// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockBlueprintSpecValidationUseCase is an autogenerated mock type for the blueprintSpecValidationUseCase type
type mockBlueprintSpecValidationUseCase struct {
	mock.Mock
}

type mockBlueprintSpecValidationUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockBlueprintSpecValidationUseCase) EXPECT() *mockBlueprintSpecValidationUseCase_Expecter {
	return &mockBlueprintSpecValidationUseCase_Expecter{mock: &_m.Mock}
}

// ValidateBlueprintSpecDynamically provides a mock function with given fields: ctx, blueprintId
func (_m *mockBlueprintSpecValidationUseCase) ValidateBlueprintSpecDynamically(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for ValidateBlueprintSpecDynamically")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateBlueprintSpecDynamically'
type mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call struct {
	*mock.Call
}

// ValidateBlueprintSpecDynamically is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockBlueprintSpecValidationUseCase_Expecter) ValidateBlueprintSpecDynamically(ctx interface{}, blueprintId interface{}) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call {
	return &mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call{Call: _e.mock.On("ValidateBlueprintSpecDynamically", ctx, blueprintId)}
}

func (_c *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call) Run(run func(ctx context.Context, blueprintId string)) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call) Return(_a0 error) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call) RunAndReturn(run func(context.Context, string) error) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecDynamically_Call {
	_c.Call.Return(run)
	return _c
}

// ValidateBlueprintSpecStatically provides a mock function with given fields: ctx, blueprintId
func (_m *mockBlueprintSpecValidationUseCase) ValidateBlueprintSpecStatically(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for ValidateBlueprintSpecStatically")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ValidateBlueprintSpecStatically'
type mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call struct {
	*mock.Call
}

// ValidateBlueprintSpecStatically is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockBlueprintSpecValidationUseCase_Expecter) ValidateBlueprintSpecStatically(ctx interface{}, blueprintId interface{}) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call {
	return &mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call{Call: _e.mock.On("ValidateBlueprintSpecStatically", ctx, blueprintId)}
}

func (_c *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call) Run(run func(ctx context.Context, blueprintId string)) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call) Return(_a0 error) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call) RunAndReturn(run func(context.Context, string) error) *mockBlueprintSpecValidationUseCase_ValidateBlueprintSpecStatically_Call {
	_c.Call.Return(run)
	return _c
}

// newMockBlueprintSpecValidationUseCase creates a new instance of mockBlueprintSpecValidationUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockBlueprintSpecValidationUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockBlueprintSpecValidationUseCase {
	mock := &mockBlueprintSpecValidationUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
