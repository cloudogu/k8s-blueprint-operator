// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	domain "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	ecosystem "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// mockComponentInstallationUseCase is an autogenerated mock type for the componentInstallationUseCase type
type mockComponentInstallationUseCase struct {
	mock.Mock
}

type mockComponentInstallationUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *mockComponentInstallationUseCase) EXPECT() *mockComponentInstallationUseCase_Expecter {
	return &mockComponentInstallationUseCase_Expecter{mock: &_m.Mock}
}

// ApplyComponentStates provides a mock function with given fields: ctx, blueprintId
func (_m *mockComponentInstallationUseCase) ApplyComponentStates(ctx context.Context, blueprintId string) error {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for ApplyComponentStates")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInstallationUseCase_ApplyComponentStates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ApplyComponentStates'
type mockComponentInstallationUseCase_ApplyComponentStates_Call struct {
	*mock.Call
}

// ApplyComponentStates is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockComponentInstallationUseCase_Expecter) ApplyComponentStates(ctx interface{}, blueprintId interface{}) *mockComponentInstallationUseCase_ApplyComponentStates_Call {
	return &mockComponentInstallationUseCase_ApplyComponentStates_Call{Call: _e.mock.On("ApplyComponentStates", ctx, blueprintId)}
}

func (_c *mockComponentInstallationUseCase_ApplyComponentStates_Call) Run(run func(ctx context.Context, blueprintId string)) *mockComponentInstallationUseCase_ApplyComponentStates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockComponentInstallationUseCase_ApplyComponentStates_Call) Return(_a0 error) *mockComponentInstallationUseCase_ApplyComponentStates_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInstallationUseCase_ApplyComponentStates_Call) RunAndReturn(run func(context.Context, string) error) *mockComponentInstallationUseCase_ApplyComponentStates_Call {
	_c.Call.Return(run)
	return _c
}

// CheckComponentHealth provides a mock function with given fields: ctx
func (_m *mockComponentInstallationUseCase) CheckComponentHealth(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CheckComponentHealth")
	}

	var r0 ecosystem.ComponentHealthResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (ecosystem.ComponentHealthResult, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) ecosystem.ComponentHealthResult); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(ecosystem.ComponentHealthResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInstallationUseCase_CheckComponentHealth_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CheckComponentHealth'
type mockComponentInstallationUseCase_CheckComponentHealth_Call struct {
	*mock.Call
}

// CheckComponentHealth is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockComponentInstallationUseCase_Expecter) CheckComponentHealth(ctx interface{}) *mockComponentInstallationUseCase_CheckComponentHealth_Call {
	return &mockComponentInstallationUseCase_CheckComponentHealth_Call{Call: _e.mock.On("CheckComponentHealth", ctx)}
}

func (_c *mockComponentInstallationUseCase_CheckComponentHealth_Call) Run(run func(ctx context.Context)) *mockComponentInstallationUseCase_CheckComponentHealth_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockComponentInstallationUseCase_CheckComponentHealth_Call) Return(_a0 ecosystem.ComponentHealthResult, _a1 error) *mockComponentInstallationUseCase_CheckComponentHealth_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInstallationUseCase_CheckComponentHealth_Call) RunAndReturn(run func(context.Context) (ecosystem.ComponentHealthResult, error)) *mockComponentInstallationUseCase_CheckComponentHealth_Call {
	_c.Call.Return(run)
	return _c
}

// WaitForHealthyComponents provides a mock function with given fields: ctx
func (_m *mockComponentInstallationUseCase) WaitForHealthyComponents(ctx context.Context) (ecosystem.ComponentHealthResult, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for WaitForHealthyComponents")
	}

	var r0 ecosystem.ComponentHealthResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (ecosystem.ComponentHealthResult, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) ecosystem.ComponentHealthResult); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(ecosystem.ComponentHealthResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockComponentInstallationUseCase_WaitForHealthyComponents_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WaitForHealthyComponents'
type mockComponentInstallationUseCase_WaitForHealthyComponents_Call struct {
	*mock.Call
}

// WaitForHealthyComponents is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockComponentInstallationUseCase_Expecter) WaitForHealthyComponents(ctx interface{}) *mockComponentInstallationUseCase_WaitForHealthyComponents_Call {
	return &mockComponentInstallationUseCase_WaitForHealthyComponents_Call{Call: _e.mock.On("WaitForHealthyComponents", ctx)}
}

func (_c *mockComponentInstallationUseCase_WaitForHealthyComponents_Call) Run(run func(ctx context.Context)) *mockComponentInstallationUseCase_WaitForHealthyComponents_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockComponentInstallationUseCase_WaitForHealthyComponents_Call) Return(_a0 ecosystem.ComponentHealthResult, _a1 error) *mockComponentInstallationUseCase_WaitForHealthyComponents_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInstallationUseCase_WaitForHealthyComponents_Call) RunAndReturn(run func(context.Context) (ecosystem.ComponentHealthResult, error)) *mockComponentInstallationUseCase_WaitForHealthyComponents_Call {
	_c.Call.Return(run)
	return _c
}

// applyComponentState provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockComponentInstallationUseCase) applyComponentState(_a0 context.Context, _a1 domain.ComponentDiff, _a2 *ecosystem.ComponentInstallation) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for applyComponentState")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.ComponentDiff, *ecosystem.ComponentInstallation) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInstallationUseCase_applyComponentState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'applyComponentState'
type mockComponentInstallationUseCase_applyComponentState_Call struct {
	*mock.Call
}

// applyComponentState is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 domain.ComponentDiff
//   - _a2 *ecosystem.ComponentInstallation
func (_e *mockComponentInstallationUseCase_Expecter) applyComponentState(_a0 interface{}, _a1 interface{}, _a2 interface{}) *mockComponentInstallationUseCase_applyComponentState_Call {
	return &mockComponentInstallationUseCase_applyComponentState_Call{Call: _e.mock.On("applyComponentState", _a0, _a1, _a2)}
}

func (_c *mockComponentInstallationUseCase_applyComponentState_Call) Run(run func(_a0 context.Context, _a1 domain.ComponentDiff, _a2 *ecosystem.ComponentInstallation)) *mockComponentInstallationUseCase_applyComponentState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.ComponentDiff), args[2].(*ecosystem.ComponentInstallation))
	})
	return _c
}

func (_c *mockComponentInstallationUseCase_applyComponentState_Call) Return(_a0 error) *mockComponentInstallationUseCase_applyComponentState_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInstallationUseCase_applyComponentState_Call) RunAndReturn(run func(context.Context, domain.ComponentDiff, *ecosystem.ComponentInstallation) error) *mockComponentInstallationUseCase_applyComponentState_Call {
	_c.Call.Return(run)
	return _c
}

// newMockComponentInstallationUseCase creates a new instance of mockComponentInstallationUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockComponentInstallationUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockComponentInstallationUseCase {
	mock := &mockComponentInstallationUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
