// Code generated by mockery v2.42.1. DO NOT EDIT.

package application

import (
	context "context"

	config "github.com/cloudogu/k8s-registry-lib/config"

	mock "github.com/stretchr/testify/mock"
)

// mockSensitiveDoguConfigRepository is an autogenerated mock type for the sensitiveDoguConfigRepository type
type mockSensitiveDoguConfigRepository struct {
	mock.Mock
}

type mockSensitiveDoguConfigRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockSensitiveDoguConfigRepository) EXPECT() *mockSensitiveDoguConfigRepository_Expecter {
	return &mockSensitiveDoguConfigRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *mockSensitiveDoguConfigRepository) Create(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
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

// mockSensitiveDoguConfigRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockSensitiveDoguConfigRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *mockSensitiveDoguConfigRepository_Expecter) Create(ctx interface{}, _a1 interface{}) *mockSensitiveDoguConfigRepository_Create_Call {
	return &mockSensitiveDoguConfigRepository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *mockSensitiveDoguConfigRepository_Create_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *mockSensitiveDoguConfigRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_Create_Call) Return(_a0 config.DoguConfig, _a1 error) *mockSensitiveDoguConfigRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_Create_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *mockSensitiveDoguConfigRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, doguName
func (_m *mockSensitiveDoguConfigRepository) Get(ctx context.Context, doguName config.SimpleDoguName) (config.DoguConfig, error) {
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

// mockSensitiveDoguConfigRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockSensitiveDoguConfigRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName config.SimpleDoguName
func (_e *mockSensitiveDoguConfigRepository_Expecter) Get(ctx interface{}, doguName interface{}) *mockSensitiveDoguConfigRepository_Get_Call {
	return &mockSensitiveDoguConfigRepository_Get_Call{Call: _e.mock.On("Get", ctx, doguName)}
}

func (_c *mockSensitiveDoguConfigRepository_Get_Call) Run(run func(ctx context.Context, doguName config.SimpleDoguName)) *mockSensitiveDoguConfigRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.SimpleDoguName))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_Get_Call) Return(_a0 config.DoguConfig, _a1 error) *mockSensitiveDoguConfigRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_Get_Call) RunAndReturn(run func(context.Context, config.SimpleDoguName) (config.DoguConfig, error)) *mockSensitiveDoguConfigRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx, doguNames
func (_m *mockSensitiveDoguConfigRepository) GetAll(ctx context.Context, doguNames []config.SimpleDoguName) (map[config.SimpleDoguName]config.DoguConfig, error) {
	ret := _m.Called(ctx, doguNames)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 map[config.SimpleDoguName]config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []config.SimpleDoguName) (map[config.SimpleDoguName]config.DoguConfig, error)); ok {
		return rf(ctx, doguNames)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []config.SimpleDoguName) map[config.SimpleDoguName]config.DoguConfig); ok {
		r0 = rf(ctx, doguNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[config.SimpleDoguName]config.DoguConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []config.SimpleDoguName) error); ok {
		r1 = rf(ctx, doguNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockSensitiveDoguConfigRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type mockSensitiveDoguConfigRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
//   - doguNames []config.SimpleDoguName
func (_e *mockSensitiveDoguConfigRepository_Expecter) GetAll(ctx interface{}, doguNames interface{}) *mockSensitiveDoguConfigRepository_GetAll_Call {
	return &mockSensitiveDoguConfigRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx, doguNames)}
}

func (_c *mockSensitiveDoguConfigRepository_GetAll_Call) Run(run func(ctx context.Context, doguNames []config.SimpleDoguName)) *mockSensitiveDoguConfigRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]config.SimpleDoguName))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_GetAll_Call) Return(_a0 map[config.SimpleDoguName]config.DoguConfig, _a1 error) *mockSensitiveDoguConfigRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_GetAll_Call) RunAndReturn(run func(context.Context, []config.SimpleDoguName) (map[config.SimpleDoguName]config.DoguConfig, error)) *mockSensitiveDoguConfigRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *mockSensitiveDoguConfigRepository) Update(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
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

// mockSensitiveDoguConfigRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockSensitiveDoguConfigRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *mockSensitiveDoguConfigRepository_Expecter) Update(ctx interface{}, _a1 interface{}) *mockSensitiveDoguConfigRepository_Update_Call {
	return &mockSensitiveDoguConfigRepository_Update_Call{Call: _e.mock.On("Update", ctx, _a1)}
}

func (_c *mockSensitiveDoguConfigRepository_Update_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *mockSensitiveDoguConfigRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_Update_Call) Return(_a0 config.DoguConfig, _a1 error) *mockSensitiveDoguConfigRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockSensitiveDoguConfigRepository_Update_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *mockSensitiveDoguConfigRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// newMockSensitiveDoguConfigRepository creates a new instance of mockSensitiveDoguConfigRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockSensitiveDoguConfigRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockSensitiveDoguConfigRepository {
	mock := &mockSensitiveDoguConfigRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
