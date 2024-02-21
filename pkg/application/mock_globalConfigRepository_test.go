// Code generated by mockery v2.20.0. DO NOT EDIT.

package application

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// mockGlobalConfigRepository is an autogenerated mock type for the globalConfigRepository type
type mockGlobalConfigRepository struct {
	mock.Mock
}

type mockGlobalConfigRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockGlobalConfigRepository) EXPECT() *mockGlobalConfigRepository_Expecter {
	return &mockGlobalConfigRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *mockGlobalConfigRepository) Delete(_a0 context.Context, _a1 common.GlobalConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.GlobalConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockGlobalConfigRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockGlobalConfigRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.GlobalConfigKey
func (_e *mockGlobalConfigRepository_Expecter) Delete(_a0 interface{}, _a1 interface{}) *mockGlobalConfigRepository_Delete_Call {
	return &mockGlobalConfigRepository_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *mockGlobalConfigRepository_Delete_Call) Run(run func(_a0 context.Context, _a1 common.GlobalConfigKey)) *mockGlobalConfigRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.GlobalConfigKey))
	})
	return _c
}

func (_c *mockGlobalConfigRepository_Delete_Call) Return(_a0 error) *mockGlobalConfigRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockGlobalConfigRepository_Delete_Call) RunAndReturn(run func(context.Context, common.GlobalConfigKey) error) *mockGlobalConfigRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAllByKeys provides a mock function with given fields: _a0, _a1
func (_m *mockGlobalConfigRepository) DeleteAllByKeys(_a0 context.Context, _a1 []common.GlobalConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.GlobalConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockGlobalConfigRepository_DeleteAllByKeys_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAllByKeys'
type mockGlobalConfigRepository_DeleteAllByKeys_Call struct {
	*mock.Call
}

// DeleteAllByKeys is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.GlobalConfigKey
func (_e *mockGlobalConfigRepository_Expecter) DeleteAllByKeys(_a0 interface{}, _a1 interface{}) *mockGlobalConfigRepository_DeleteAllByKeys_Call {
	return &mockGlobalConfigRepository_DeleteAllByKeys_Call{Call: _e.mock.On("DeleteAllByKeys", _a0, _a1)}
}

func (_c *mockGlobalConfigRepository_DeleteAllByKeys_Call) Run(run func(_a0 context.Context, _a1 []common.GlobalConfigKey)) *mockGlobalConfigRepository_DeleteAllByKeys_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.GlobalConfigKey))
	})
	return _c
}

func (_c *mockGlobalConfigRepository_DeleteAllByKeys_Call) Return(_a0 error) *mockGlobalConfigRepository_DeleteAllByKeys_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockGlobalConfigRepository_DeleteAllByKeys_Call) RunAndReturn(run func(context.Context, []common.GlobalConfigKey) error) *mockGlobalConfigRepository_DeleteAllByKeys_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *mockGlobalConfigRepository) Get(_a0 context.Context, _a1 common.GlobalConfigKey) (*ecosystem.GlobalConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ecosystem.GlobalConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.GlobalConfigKey) (*ecosystem.GlobalConfigEntry, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.GlobalConfigKey) *ecosystem.GlobalConfigEntry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.GlobalConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.GlobalConfigKey) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockGlobalConfigRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockGlobalConfigRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.GlobalConfigKey
func (_e *mockGlobalConfigRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *mockGlobalConfigRepository_Get_Call {
	return &mockGlobalConfigRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *mockGlobalConfigRepository_Get_Call) Run(run func(_a0 context.Context, _a1 common.GlobalConfigKey)) *mockGlobalConfigRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.GlobalConfigKey))
	})
	return _c
}

func (_c *mockGlobalConfigRepository_Get_Call) Return(_a0 *ecosystem.GlobalConfigEntry, _a1 error) *mockGlobalConfigRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockGlobalConfigRepository_Get_Call) RunAndReturn(run func(context.Context, common.GlobalConfigKey) (*ecosystem.GlobalConfigEntry, error)) *mockGlobalConfigRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByKey provides a mock function with given fields: _a0, _a1
func (_m *mockGlobalConfigRepository) GetAllByKey(_a0 context.Context, _a1 []common.GlobalConfigKey) (map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

	var r0 map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.GlobalConfigKey) (map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []common.GlobalConfigKey) map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []common.GlobalConfigKey) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockGlobalConfigRepository_GetAllByKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByKey'
type mockGlobalConfigRepository_GetAllByKey_Call struct {
	*mock.Call
}

// GetAllByKey is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.GlobalConfigKey
func (_e *mockGlobalConfigRepository_Expecter) GetAllByKey(_a0 interface{}, _a1 interface{}) *mockGlobalConfigRepository_GetAllByKey_Call {
	return &mockGlobalConfigRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", _a0, _a1)}
}

func (_c *mockGlobalConfigRepository_GetAllByKey_Call) Run(run func(_a0 context.Context, _a1 []common.GlobalConfigKey)) *mockGlobalConfigRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.GlobalConfigKey))
	})
	return _c
}

func (_c *mockGlobalConfigRepository_GetAllByKey_Call) Return(_a0 map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry, _a1 error) *mockGlobalConfigRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockGlobalConfigRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []common.GlobalConfigKey) (map[common.GlobalConfigKey]*ecosystem.GlobalConfigEntry, error)) *mockGlobalConfigRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *mockGlobalConfigRepository) Save(_a0 context.Context, _a1 *ecosystem.GlobalConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.GlobalConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockGlobalConfigRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type mockGlobalConfigRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ecosystem.GlobalConfigEntry
func (_e *mockGlobalConfigRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *mockGlobalConfigRepository_Save_Call {
	return &mockGlobalConfigRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *mockGlobalConfigRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *ecosystem.GlobalConfigEntry)) *mockGlobalConfigRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.GlobalConfigEntry))
	})
	return _c
}

func (_c *mockGlobalConfigRepository_Save_Call) Return(_a0 error) *mockGlobalConfigRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockGlobalConfigRepository_Save_Call) RunAndReturn(run func(context.Context, *ecosystem.GlobalConfigEntry) error) *mockGlobalConfigRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: _a0, _a1
func (_m *mockGlobalConfigRepository) SaveAll(_a0 context.Context, _a1 []*ecosystem.GlobalConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.GlobalConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockGlobalConfigRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type mockGlobalConfigRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []*ecosystem.GlobalConfigEntry
func (_e *mockGlobalConfigRepository_Expecter) SaveAll(_a0 interface{}, _a1 interface{}) *mockGlobalConfigRepository_SaveAll_Call {
	return &mockGlobalConfigRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", _a0, _a1)}
}

func (_c *mockGlobalConfigRepository_SaveAll_Call) Run(run func(_a0 context.Context, _a1 []*ecosystem.GlobalConfigEntry)) *mockGlobalConfigRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*ecosystem.GlobalConfigEntry))
	})
	return _c
}

func (_c *mockGlobalConfigRepository_SaveAll_Call) Return(_a0 error) *mockGlobalConfigRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockGlobalConfigRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []*ecosystem.GlobalConfigEntry) error) *mockGlobalConfigRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockGlobalConfigRepository interface {
	mock.TestingT
	Cleanup(func())
}

// newMockGlobalConfigRepository creates a new instance of mockGlobalConfigRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockGlobalConfigRepository(t mockConstructorTestingTnewMockGlobalConfigRepository) *mockGlobalConfigRepository {
	mock := &mockGlobalConfigRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
