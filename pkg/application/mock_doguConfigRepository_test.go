// Code generated by mockery v2.42.1. DO NOT EDIT.

package application

import (
	context "context"

	config "github.com/cloudogu/k8s-registry-lib/config"

	dogu "github.com/cloudogu/ces-commons-lib/dogu"

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

// Create provides a mock function with given fields: ctx, _a1
func (_m *mockDoguConfigRepository) Create(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
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

// mockDoguConfigRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockDoguConfigRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *mockDoguConfigRepository_Expecter) Create(ctx interface{}, _a1 interface{}) *mockDoguConfigRepository_Create_Call {
	return &mockDoguConfigRepository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *mockDoguConfigRepository_Create_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *mockDoguConfigRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *mockDoguConfigRepository_Create_Call) Return(_a0 config.DoguConfig, _a1 error) *mockDoguConfigRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_Create_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *mockDoguConfigRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, doguName
func (_m *mockDoguConfigRepository) Get(ctx context.Context, doguName dogu.SimpleName) (config.DoguConfig, error) {
	ret := _m.Called(ctx, doguName)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleName) (config.DoguConfig, error)); ok {
		return rf(ctx, doguName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleName) config.DoguConfig); ok {
		r0 = rf(ctx, doguName)
	} else {
		r0 = ret.Get(0).(config.DoguConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dogu.SimpleName) error); ok {
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
//   - doguName dogu.SimpleName
func (_e *mockDoguConfigRepository_Expecter) Get(ctx interface{}, doguName interface{}) *mockDoguConfigRepository_Get_Call {
	return &mockDoguConfigRepository_Get_Call{Call: _e.mock.On("Get", ctx, doguName)}
}

func (_c *mockDoguConfigRepository_Get_Call) Run(run func(ctx context.Context, doguName dogu.SimpleName)) *mockDoguConfigRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.SimpleName))
	})
	return _c
}

func (_c *mockDoguConfigRepository_Get_Call) Return(_a0 config.DoguConfig, _a1 error) *mockDoguConfigRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_Get_Call) RunAndReturn(run func(context.Context, dogu.SimpleName) (config.DoguConfig, error)) *mockDoguConfigRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx, doguNames
func (_m *mockDoguConfigRepository) GetAll(ctx context.Context, doguNames []dogu.SimpleName) (map[dogu.SimpleName]config.DoguConfig, error) {
	ret := _m.Called(ctx, doguNames)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 map[dogu.SimpleName]config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleName) (map[dogu.SimpleName]config.DoguConfig, error)); ok {
		return rf(ctx, doguNames)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleName) map[dogu.SimpleName]config.DoguConfig); ok {
		r0 = rf(ctx, doguNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[dogu.SimpleName]config.DoguConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []dogu.SimpleName) error); ok {
		r1 = rf(ctx, doguNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type mockDoguConfigRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
//   - doguNames []dogu.SimpleName
func (_e *mockDoguConfigRepository_Expecter) GetAll(ctx interface{}, doguNames interface{}) *mockDoguConfigRepository_GetAll_Call {
	return &mockDoguConfigRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx, doguNames)}
}

func (_c *mockDoguConfigRepository_GetAll_Call) Run(run func(ctx context.Context, doguNames []dogu.SimpleName)) *mockDoguConfigRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]dogu.SimpleName))
	})
	return _c
}

func (_c *mockDoguConfigRepository_GetAll_Call) Return(_a0 map[dogu.SimpleName]config.DoguConfig, _a1 error) *mockDoguConfigRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_GetAll_Call) RunAndReturn(run func(context.Context, []dogu.SimpleName) (map[dogu.SimpleName]config.DoguConfig, error)) *mockDoguConfigRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllExisting provides a mock function with given fields: ctx, doguNames
func (_m *mockDoguConfigRepository) GetAllExisting(ctx context.Context, doguNames []dogu.SimpleName) (map[dogu.SimpleName]config.DoguConfig, error) {
	ret := _m.Called(ctx, doguNames)

	if len(ret) == 0 {
		panic("no return value specified for GetAllExisting")
	}

	var r0 map[dogu.SimpleName]config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleName) (map[dogu.SimpleName]config.DoguConfig, error)); ok {
		return rf(ctx, doguNames)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleName) map[dogu.SimpleName]config.DoguConfig); ok {
		r0 = rf(ctx, doguNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[dogu.SimpleName]config.DoguConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []dogu.SimpleName) error); ok {
		r1 = rf(ctx, doguNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigRepository_GetAllExisting_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllExisting'
type mockDoguConfigRepository_GetAllExisting_Call struct {
	*mock.Call
}

// GetAllExisting is a helper method to define mock.On call
//   - ctx context.Context
//   - doguNames []dogu.SimpleName
func (_e *mockDoguConfigRepository_Expecter) GetAllExisting(ctx interface{}, doguNames interface{}) *mockDoguConfigRepository_GetAllExisting_Call {
	return &mockDoguConfigRepository_GetAllExisting_Call{Call: _e.mock.On("GetAllExisting", ctx, doguNames)}
}

func (_c *mockDoguConfigRepository_GetAllExisting_Call) Run(run func(ctx context.Context, doguNames []dogu.SimpleName)) *mockDoguConfigRepository_GetAllExisting_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]dogu.SimpleName))
	})
	return _c
}

func (_c *mockDoguConfigRepository_GetAllExisting_Call) Return(_a0 map[dogu.SimpleName]config.DoguConfig, _a1 error) *mockDoguConfigRepository_GetAllExisting_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_GetAllExisting_Call) RunAndReturn(run func(context.Context, []dogu.SimpleName) (map[dogu.SimpleName]config.DoguConfig, error)) *mockDoguConfigRepository_GetAllExisting_Call {
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

// UpdateOrCreate provides a mock function with given fields: ctx, _a1
func (_m *mockDoguConfigRepository) UpdateOrCreate(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrCreate")
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

// mockDoguConfigRepository_UpdateOrCreate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOrCreate'
type mockDoguConfigRepository_UpdateOrCreate_Call struct {
	*mock.Call
}

// UpdateOrCreate is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *mockDoguConfigRepository_Expecter) UpdateOrCreate(ctx interface{}, _a1 interface{}) *mockDoguConfigRepository_UpdateOrCreate_Call {
	return &mockDoguConfigRepository_UpdateOrCreate_Call{Call: _e.mock.On("UpdateOrCreate", ctx, _a1)}
}

func (_c *mockDoguConfigRepository_UpdateOrCreate_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *mockDoguConfigRepository_UpdateOrCreate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *mockDoguConfigRepository_UpdateOrCreate_Call) Return(_a0 config.DoguConfig, _a1 error) *mockDoguConfigRepository_UpdateOrCreate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigRepository_UpdateOrCreate_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *mockDoguConfigRepository_UpdateOrCreate_Call {
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
