// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	dogu "github.com/cloudogu/ces-commons-lib/dogu"
	mock "github.com/stretchr/testify/mock"
)

// mockDoguRestartRepository is an autogenerated mock type for the doguRestartRepository type
type mockDoguRestartRepository struct {
	mock.Mock
}

type mockDoguRestartRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguRestartRepository) EXPECT() *mockDoguRestartRepository_Expecter {
	return &mockDoguRestartRepository_Expecter{mock: &_m.Mock}
}

// RestartAll provides a mock function with given fields: _a0, _a1
func (_m *mockDoguRestartRepository) RestartAll(_a0 context.Context, _a1 []dogu.SimpleName) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for RestartAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleName) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguRestartRepository_RestartAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RestartAll'
type mockDoguRestartRepository_RestartAll_Call struct {
	*mock.Call
}

// RestartAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []dogu.SimpleName
func (_e *mockDoguRestartRepository_Expecter) RestartAll(_a0 interface{}, _a1 interface{}) *mockDoguRestartRepository_RestartAll_Call {
	return &mockDoguRestartRepository_RestartAll_Call{Call: _e.mock.On("RestartAll", _a0, _a1)}
}

func (_c *mockDoguRestartRepository_RestartAll_Call) Run(run func(_a0 context.Context, _a1 []dogu.SimpleName)) *mockDoguRestartRepository_RestartAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]dogu.SimpleName))
	})
	return _c
}

func (_c *mockDoguRestartRepository_RestartAll_Call) Return(_a0 error) *mockDoguRestartRepository_RestartAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguRestartRepository_RestartAll_Call) RunAndReturn(run func(context.Context, []dogu.SimpleName) error) *mockDoguRestartRepository_RestartAll_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguRestartRepository creates a new instance of mockDoguRestartRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguRestartRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguRestartRepository {
	mock := &mockDoguRestartRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
