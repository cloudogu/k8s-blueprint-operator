// Code generated by mockery v2.20.0. DO NOT EDIT.

package application

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// mockDoguConfigEntryRepository is an autogenerated mock type for the doguConfigEntryRepository type
type mockDoguConfigEntryRepository struct {
	mock.Mock
}

type mockDoguConfigEntryRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguConfigEntryRepository) EXPECT() *mockDoguConfigEntryRepository_Expecter {
	return &mockDoguConfigEntryRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *mockDoguConfigEntryRepository) Delete(_a0 context.Context, _a1 common.DoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.DoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigEntryRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDoguConfigEntryRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.DoguConfigKey
func (_e *mockDoguConfigEntryRepository_Expecter) Delete(_a0 interface{}, _a1 interface{}) *mockDoguConfigEntryRepository_Delete_Call {
	return &mockDoguConfigEntryRepository_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *mockDoguConfigEntryRepository_Delete_Call) Run(run func(_a0 context.Context, _a1 common.DoguConfigKey)) *mockDoguConfigEntryRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.DoguConfigKey))
	})
	return _c
}

func (_c *mockDoguConfigEntryRepository_Delete_Call) Return(_a0 error) *mockDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigEntryRepository_Delete_Call) RunAndReturn(run func(context.Context, common.DoguConfigKey) error) *mockDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAllByKeys provides a mock function with given fields: _a0, _a1
func (_m *mockDoguConfigEntryRepository) DeleteAllByKeys(_a0 context.Context, _a1 []common.DoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.DoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigEntryRepository_DeleteAllByKeys_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAllByKeys'
type mockDoguConfigEntryRepository_DeleteAllByKeys_Call struct {
	*mock.Call
}

// DeleteAllByKeys is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.DoguConfigKey
func (_e *mockDoguConfigEntryRepository_Expecter) DeleteAllByKeys(_a0 interface{}, _a1 interface{}) *mockDoguConfigEntryRepository_DeleteAllByKeys_Call {
	return &mockDoguConfigEntryRepository_DeleteAllByKeys_Call{Call: _e.mock.On("DeleteAllByKeys", _a0, _a1)}
}

func (_c *mockDoguConfigEntryRepository_DeleteAllByKeys_Call) Run(run func(_a0 context.Context, _a1 []common.DoguConfigKey)) *mockDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.DoguConfigKey))
	})
	return _c
}

func (_c *mockDoguConfigEntryRepository_DeleteAllByKeys_Call) Return(_a0 error) *mockDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigEntryRepository_DeleteAllByKeys_Call) RunAndReturn(run func(context.Context, []common.DoguConfigKey) error) *mockDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *mockDoguConfigEntryRepository) Get(_a0 context.Context, _a1 common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *ecosystem.DoguConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.DoguConfigKey) *ecosystem.DoguConfigEntry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.DoguConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.DoguConfigKey) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigEntryRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockDoguConfigEntryRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.DoguConfigKey
func (_e *mockDoguConfigEntryRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *mockDoguConfigEntryRepository_Get_Call {
	return &mockDoguConfigEntryRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *mockDoguConfigEntryRepository_Get_Call) Run(run func(_a0 context.Context, _a1 common.DoguConfigKey)) *mockDoguConfigEntryRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.DoguConfigKey))
	})
	return _c
}

func (_c *mockDoguConfigEntryRepository_Get_Call) Return(_a0 *ecosystem.DoguConfigEntry, _a1 error) *mockDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigEntryRepository_Get_Call) RunAndReturn(run func(context.Context, common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error)) *mockDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByKey provides a mock function with given fields: _a0, _a1
func (_m *mockDoguConfigEntryRepository) GetAllByKey(_a0 context.Context, _a1 []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

	var r0 map[common.DoguConfigKey]*ecosystem.DoguConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []common.DoguConfigKey) map[common.DoguConfigKey]*ecosystem.DoguConfigEntry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.DoguConfigKey]*ecosystem.DoguConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []common.DoguConfigKey) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguConfigEntryRepository_GetAllByKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByKey'
type mockDoguConfigEntryRepository_GetAllByKey_Call struct {
	*mock.Call
}

// GetAllByKey is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.DoguConfigKey
func (_e *mockDoguConfigEntryRepository_Expecter) GetAllByKey(_a0 interface{}, _a1 interface{}) *mockDoguConfigEntryRepository_GetAllByKey_Call {
	return &mockDoguConfigEntryRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", _a0, _a1)}
}

func (_c *mockDoguConfigEntryRepository_GetAllByKey_Call) Run(run func(_a0 context.Context, _a1 []common.DoguConfigKey)) *mockDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.DoguConfigKey))
	})
	return _c
}

func (_c *mockDoguConfigEntryRepository_GetAllByKey_Call) Return(_a0 map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, _a1 error) *mockDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguConfigEntryRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error)) *mockDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *mockDoguConfigEntryRepository) Save(_a0 context.Context, _a1 *ecosystem.DoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.DoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigEntryRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type mockDoguConfigEntryRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ecosystem.DoguConfigEntry
func (_e *mockDoguConfigEntryRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *mockDoguConfigEntryRepository_Save_Call {
	return &mockDoguConfigEntryRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *mockDoguConfigEntryRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *ecosystem.DoguConfigEntry)) *mockDoguConfigEntryRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.DoguConfigEntry))
	})
	return _c
}

func (_c *mockDoguConfigEntryRepository_Save_Call) Return(_a0 error) *mockDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigEntryRepository_Save_Call) RunAndReturn(run func(context.Context, *ecosystem.DoguConfigEntry) error) *mockDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: _a0, _a1
func (_m *mockDoguConfigEntryRepository) SaveAll(_a0 context.Context, _a1 []*ecosystem.DoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.DoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguConfigEntryRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type mockDoguConfigEntryRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []*ecosystem.DoguConfigEntry
func (_e *mockDoguConfigEntryRepository_Expecter) SaveAll(_a0 interface{}, _a1 interface{}) *mockDoguConfigEntryRepository_SaveAll_Call {
	return &mockDoguConfigEntryRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", _a0, _a1)}
}

func (_c *mockDoguConfigEntryRepository_SaveAll_Call) Run(run func(_a0 context.Context, _a1 []*ecosystem.DoguConfigEntry)) *mockDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*ecosystem.DoguConfigEntry))
	})
	return _c
}

func (_c *mockDoguConfigEntryRepository_SaveAll_Call) Return(_a0 error) *mockDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguConfigEntryRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []*ecosystem.DoguConfigEntry) error) *mockDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguConfigEntryRepository interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguConfigEntryRepository creates a new instance of mockDoguConfigEntryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguConfigEntryRepository(t mockConstructorTestingTnewMockDoguConfigEntryRepository) *mockDoguConfigEntryRepository {
	mock := &mockDoguConfigEntryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
