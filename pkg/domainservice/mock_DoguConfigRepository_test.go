// Code generated by mockery v2.42.1. DO NOT EDIT.

package domainservice

import (
	context "context"

	config "github.com/cloudogu/k8s-registry-lib/config"

	dogu "github.com/cloudogu/ces-commons-lib/dogu"

	mock "github.com/stretchr/testify/mock"
)

// MockDoguConfigRepository is an autogenerated mock type for the DoguConfigRepository type
type MockDoguConfigRepository struct {
	mock.Mock
}

type MockDoguConfigRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDoguConfigRepository) EXPECT() *MockDoguConfigRepository_Expecter {
	return &MockDoguConfigRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *MockDoguConfigRepository) Create(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
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

// MockDoguConfigRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockDoguConfigRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *MockDoguConfigRepository_Expecter) Create(ctx interface{}, _a1 interface{}) *MockDoguConfigRepository_Create_Call {
	return &MockDoguConfigRepository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *MockDoguConfigRepository_Create_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *MockDoguConfigRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *MockDoguConfigRepository_Create_Call) Return(_a0 config.DoguConfig, _a1 error) *MockDoguConfigRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigRepository_Create_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *MockDoguConfigRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, doguName
func (_m *MockDoguConfigRepository) Get(ctx context.Context, doguName dogu.SimpleDoguName) (config.DoguConfig, error) {
	ret := _m.Called(ctx, doguName)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleDoguName) (config.DoguConfig, error)); ok {
		return rf(ctx, doguName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleDoguName) config.DoguConfig); ok {
		r0 = rf(ctx, doguName)
	} else {
		r0 = ret.Get(0).(config.DoguConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dogu.SimpleDoguName) error); ok {
		r1 = rf(ctx, doguName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguConfigRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockDoguConfigRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName dogu.SimpleDoguName
func (_e *MockDoguConfigRepository_Expecter) Get(ctx interface{}, doguName interface{}) *MockDoguConfigRepository_Get_Call {
	return &MockDoguConfigRepository_Get_Call{Call: _e.mock.On("Get", ctx, doguName)}
}

func (_c *MockDoguConfigRepository_Get_Call) Run(run func(ctx context.Context, doguName dogu.SimpleDoguName)) *MockDoguConfigRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.SimpleDoguName))
	})
	return _c
}

func (_c *MockDoguConfigRepository_Get_Call) Return(_a0 config.DoguConfig, _a1 error) *MockDoguConfigRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigRepository_Get_Call) RunAndReturn(run func(context.Context, dogu.SimpleDoguName) (config.DoguConfig, error)) *MockDoguConfigRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx, doguNames
func (_m *MockDoguConfigRepository) GetAll(ctx context.Context, doguNames []dogu.SimpleDoguName) (map[dogu.SimpleDoguName]config.DoguConfig, error) {
	ret := _m.Called(ctx, doguNames)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 map[dogu.SimpleDoguName]config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleDoguName) (map[dogu.SimpleDoguName]config.DoguConfig, error)); ok {
		return rf(ctx, doguNames)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleDoguName) map[dogu.SimpleDoguName]config.DoguConfig); ok {
		r0 = rf(ctx, doguNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[dogu.SimpleDoguName]config.DoguConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []dogu.SimpleDoguName) error); ok {
		r1 = rf(ctx, doguNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguConfigRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockDoguConfigRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
//   - doguNames []dogu.SimpleDoguName
func (_e *MockDoguConfigRepository_Expecter) GetAll(ctx interface{}, doguNames interface{}) *MockDoguConfigRepository_GetAll_Call {
	return &MockDoguConfigRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx, doguNames)}
}

func (_c *MockDoguConfigRepository_GetAll_Call) Run(run func(ctx context.Context, doguNames []dogu.SimpleDoguName)) *MockDoguConfigRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]dogu.SimpleDoguName))
	})
	return _c
}

func (_c *MockDoguConfigRepository_GetAll_Call) Return(_a0 map[dogu.SimpleDoguName]config.DoguConfig, _a1 error) *MockDoguConfigRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigRepository_GetAll_Call) RunAndReturn(run func(context.Context, []dogu.SimpleDoguName) (map[dogu.SimpleDoguName]config.DoguConfig, error)) *MockDoguConfigRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllExisting provides a mock function with given fields: ctx, doguNames
func (_m *MockDoguConfigRepository) GetAllExisting(ctx context.Context, doguNames []dogu.SimpleDoguName) (map[dogu.SimpleDoguName]config.DoguConfig, error) {
	ret := _m.Called(ctx, doguNames)

	if len(ret) == 0 {
		panic("no return value specified for GetAllExisting")
	}

	var r0 map[dogu.SimpleDoguName]config.DoguConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleDoguName) (map[dogu.SimpleDoguName]config.DoguConfig, error)); ok {
		return rf(ctx, doguNames)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.SimpleDoguName) map[dogu.SimpleDoguName]config.DoguConfig); ok {
		r0 = rf(ctx, doguNames)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[dogu.SimpleDoguName]config.DoguConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []dogu.SimpleDoguName) error); ok {
		r1 = rf(ctx, doguNames)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguConfigRepository_GetAllExisting_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllExisting'
type MockDoguConfigRepository_GetAllExisting_Call struct {
	*mock.Call
}

// GetAllExisting is a helper method to define mock.On call
//   - ctx context.Context
//   - doguNames []dogu.SimpleDoguName
func (_e *MockDoguConfigRepository_Expecter) GetAllExisting(ctx interface{}, doguNames interface{}) *MockDoguConfigRepository_GetAllExisting_Call {
	return &MockDoguConfigRepository_GetAllExisting_Call{Call: _e.mock.On("GetAllExisting", ctx, doguNames)}
}

func (_c *MockDoguConfigRepository_GetAllExisting_Call) Run(run func(ctx context.Context, doguNames []dogu.SimpleDoguName)) *MockDoguConfigRepository_GetAllExisting_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]dogu.SimpleDoguName))
	})
	return _c
}

func (_c *MockDoguConfigRepository_GetAllExisting_Call) Return(_a0 map[dogu.SimpleDoguName]config.DoguConfig, _a1 error) *MockDoguConfigRepository_GetAllExisting_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigRepository_GetAllExisting_Call) RunAndReturn(run func(context.Context, []dogu.SimpleDoguName) (map[dogu.SimpleDoguName]config.DoguConfig, error)) *MockDoguConfigRepository_GetAllExisting_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *MockDoguConfigRepository) Update(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
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

// MockDoguConfigRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockDoguConfigRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *MockDoguConfigRepository_Expecter) Update(ctx interface{}, _a1 interface{}) *MockDoguConfigRepository_Update_Call {
	return &MockDoguConfigRepository_Update_Call{Call: _e.mock.On("Update", ctx, _a1)}
}

func (_c *MockDoguConfigRepository_Update_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *MockDoguConfigRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *MockDoguConfigRepository_Update_Call) Return(_a0 config.DoguConfig, _a1 error) *MockDoguConfigRepository_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigRepository_Update_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *MockDoguConfigRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOrCreate provides a mock function with given fields: ctx, _a1
func (_m *MockDoguConfigRepository) UpdateOrCreate(ctx context.Context, _a1 config.DoguConfig) (config.DoguConfig, error) {
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

// MockDoguConfigRepository_UpdateOrCreate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOrCreate'
type MockDoguConfigRepository_UpdateOrCreate_Call struct {
	*mock.Call
}

// UpdateOrCreate is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 config.DoguConfig
func (_e *MockDoguConfigRepository_Expecter) UpdateOrCreate(ctx interface{}, _a1 interface{}) *MockDoguConfigRepository_UpdateOrCreate_Call {
	return &MockDoguConfigRepository_UpdateOrCreate_Call{Call: _e.mock.On("UpdateOrCreate", ctx, _a1)}
}

func (_c *MockDoguConfigRepository_UpdateOrCreate_Call) Run(run func(ctx context.Context, _a1 config.DoguConfig)) *MockDoguConfigRepository_UpdateOrCreate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.DoguConfig))
	})
	return _c
}

func (_c *MockDoguConfigRepository_UpdateOrCreate_Call) Return(_a0 config.DoguConfig, _a1 error) *MockDoguConfigRepository_UpdateOrCreate_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigRepository_UpdateOrCreate_Call) RunAndReturn(run func(context.Context, config.DoguConfig) (config.DoguConfig, error)) *MockDoguConfigRepository_UpdateOrCreate_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDoguConfigRepository creates a new instance of MockDoguConfigRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDoguConfigRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDoguConfigRepository {
	mock := &MockDoguConfigRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
