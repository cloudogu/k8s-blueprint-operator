// Code generated by mockery v2.42.1. DO NOT EDIT.

package domainservice

import (
	context "context"

	config "github.com/cloudogu/k8s-registry-lib/config"

	mock "github.com/stretchr/testify/mock"
)

// MockDoguRestartRepository is an autogenerated mock type for the DoguRestartRepository type
type MockDoguRestartRepository struct {
	mock.Mock
}

type MockDoguRestartRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDoguRestartRepository) EXPECT() *MockDoguRestartRepository_Expecter {
	return &MockDoguRestartRepository_Expecter{mock: &_m.Mock}
}

// RestartAll provides a mock function with given fields: _a0, _a1
func (_m *MockDoguRestartRepository) RestartAll(_a0 context.Context, _a1 []config.SimpleDoguName) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for RestartAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []config.SimpleDoguName) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguRestartRepository_RestartAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RestartAll'
type MockDoguRestartRepository_RestartAll_Call struct {
	*mock.Call
}

// RestartAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []config.SimpleDoguName
func (_e *MockDoguRestartRepository_Expecter) RestartAll(_a0 interface{}, _a1 interface{}) *MockDoguRestartRepository_RestartAll_Call {
	return &MockDoguRestartRepository_RestartAll_Call{Call: _e.mock.On("RestartAll", _a0, _a1)}
}

func (_c *MockDoguRestartRepository_RestartAll_Call) Run(run func(_a0 context.Context, _a1 []config.SimpleDoguName)) *MockDoguRestartRepository_RestartAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]config.SimpleDoguName))
	})
	return _c
}

func (_c *MockDoguRestartRepository_RestartAll_Call) Return(_a0 error) *MockDoguRestartRepository_RestartAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguRestartRepository_RestartAll_Call) RunAndReturn(run func(context.Context, []config.SimpleDoguName) error) *MockDoguRestartRepository_RestartAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDoguRestartRepository creates a new instance of MockDoguRestartRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDoguRestartRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDoguRestartRepository {
	mock := &MockDoguRestartRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}