// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	mock "github.com/stretchr/testify/mock"
)

// MockDoguRestartAdapter is an autogenerated mock type for the DoguRestartRepository type
type MockDoguRestartAdapter struct {
	mock.Mock
}

type MockDoguRestartAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDoguRestartAdapter) EXPECT() *MockDoguRestartAdapter_Expecter {
	return &MockDoguRestartAdapter_Expecter{mock: &_m.Mock}
}

// RestartAll provides a mock function with given fields: _a0, _a1
func (_m *MockDoguRestartAdapter) RestartAll(_a0 context.Context, _a1 []common.SimpleDoguName) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.SimpleDoguName) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguRestartAdapter_RestartAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RestartAll'
type MockDoguRestartAdapter_RestartAll_Call struct {
	*mock.Call
}

// RestartAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.SimpleDoguName
func (_e *MockDoguRestartAdapter_Expecter) RestartAll(_a0 interface{}, _a1 interface{}) *MockDoguRestartAdapter_RestartAll_Call {
	return &MockDoguRestartAdapter_RestartAll_Call{Call: _e.mock.On("RestartAll", _a0, _a1)}
}

func (_c *MockDoguRestartAdapter_RestartAll_Call) Run(run func(_a0 context.Context, _a1 []common.SimpleDoguName)) *MockDoguRestartAdapter_RestartAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.SimpleDoguName))
	})
	return _c
}

func (_c *MockDoguRestartAdapter_RestartAll_Call) Return(_a0 error) *MockDoguRestartAdapter_RestartAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguRestartAdapter_RestartAll_Call) RunAndReturn(run func(context.Context, []common.SimpleDoguName) error) *MockDoguRestartAdapter_RestartAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDoguRestartAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDoguRestartAdapter creates a new instance of MockDoguRestartAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDoguRestartAdapter(t mockConstructorTestingTNewMockDoguRestartAdapter) *MockDoguRestartAdapter {
	mock := &MockDoguRestartAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
