// Code generated by mockery v2.42.1. DO NOT EDIT.

package application

import (
	context "context"

	config "github.com/cloudogu/k8s-registry-lib/config"

	mock "github.com/stretchr/testify/mock"
)

// mockDoguConfigRepository is an autogenerated mock type for the doguConfigRepository type
type mockDoguConfigRepository struct {
	mock.Mock
}

type mockDoguConfigRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguConfigRepository) EXPECT() *mockDoguConfigRepository_Expecter {
	return &mockDoguConfigRepository_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, doguName
func (_m *mockDoguConfigRepository) Get(ctx context.Context, doguName config.SimpleDoguName) (config.DoguConfig, error) {
	ret := _m.Called(ctx, doguName)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, config.SimpleDoguName) (config.DoguConfig, error)); ok {
		return rf(ctx, doguName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, config.SimpleDoguName) config.DoguConfig); ok {
		r0 = rf(ctx, doguName)
	} else {
		r0 = ret.Get(0).(config.DoguConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, config.SimpleDoguName) error); ok {
		r1 = rf(ctx, doguName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDoguConfigRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName config.SimpleDoguName
func (_e *mockDoguConfigRepository_Expecter) Get(ctx interface{}, doguName interface{}) *mockDoguConfigRepository_Get_Call {
	return &mockDoguConfigRepository_Get_Call{Call: _e.mock.On("Get", ctx, doguName)}
}

func (_c *mockDoguConfigRepository_Get_Call) Run(run func(ctx context.Context, doguName config.SimpleDoguName)) *mockDoguConfigRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.SimpleDoguName))
	})
	return _c
}

func (_c *mockDoguConfigRepository_Get_Call) Return(_a0 config.DoguConfig, _a1 error) *mockDoguConfigRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_Get_Call) RunAndReturn(run func(context.Context, config.SimpleDoguName) (config.DoguConfig, error)) *mockDoguConfigRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *mockDoguConfigRepository) Update(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, config.DoguConfig) (config.DoguConfig, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, config.DoguConfig) config.DoguConfig); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(config.DoguConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, config.DoguConfig) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockDoguConfigRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *mockDoguConfigRepository_Expecter) Update(ctx interface{}, _a1 interface{}) *mockDoguConfigRepository_Update_Call {
	return &mockDoguConfigRepository_Update_Call{Call: _e.mock.On("Update", ctx, _a1)}
}

func (_c *mockDoguConfigRepository_Update_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *mockDoguConfigRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *mockDoguConfigRepository_Update_Call) Return(_a0 config.DoguConfig, _a1 error) *mockDoguConfigRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_Update_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *mockDoguConfigRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguConfigRepository creates a new instance of mockDoguConfigRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguConfigRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguConfigRepository {
	mock := &mockDoguConfigRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}