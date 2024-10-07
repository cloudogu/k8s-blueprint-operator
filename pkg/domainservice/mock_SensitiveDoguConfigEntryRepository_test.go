// Code generated by mockery v2.42.1. DO NOT EDIT.

package domainservice

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// MockSensitiveDoguConfigEntryRepository is an autogenerated mock type for the SensitiveDoguConfigEntryRepository type
type MockSensitiveDoguConfigEntryRepository struct {
	mock.Mock
}

type MockSensitiveDoguConfigEntryRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSensitiveDoguConfigEntryRepository) EXPECT() *MockSensitiveDoguConfigEntryRepository_Expecter {
	return &MockSensitiveDoguConfigEntryRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) Delete(_a0 context.Context, _a1 common.SensitiveDoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SensitiveDoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSensitiveDoguConfigEntryRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockSensitiveDoguConfigEntryRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SensitiveDoguConfigKey
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) Delete(_a0 interface{}, _a1 interface{}) *MockSensitiveDoguConfigEntryRepository_Delete_Call {
	return &MockSensitiveDoguConfigEntryRepository_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_Delete_Call) Run(run func(_a0 context.Context, _a1 common.SensitiveDoguConfigKey)) *MockSensitiveDoguConfigEntryRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_Delete_Call) Return(_a0 error) *MockSensitiveDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_Delete_Call) RunAndReturn(run func(context.Context, common.SensitiveDoguConfigKey) error) *MockSensitiveDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAllByKeys provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) DeleteAllByKeys(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllByKeys")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.SensitiveDoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAllByKeys'
type MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call struct {
	*mock.Call
}

// DeleteAllByKeys is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.SensitiveDoguConfigKey
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) DeleteAllByKeys(_a0 interface{}, _a1 interface{}) *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	return &MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call{Call: _e.mock.On("DeleteAllByKeys", _a0, _a1)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call) Run(run func(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey)) *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call) Return(_a0 error) *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call) RunAndReturn(run func(context.Context, []common.SensitiveDoguConfigKey) error) *MockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) Get(_a0 context.Context, _a1 common.SensitiveDoguConfigKey) (*ecosystem.SensitiveDoguConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *ecosystem.SensitiveDoguConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SensitiveDoguConfigKey) (*ecosystem.SensitiveDoguConfigEntry, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.SensitiveDoguConfigKey) *ecosystem.SensitiveDoguConfigEntry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.SensitiveDoguConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.SensitiveDoguConfigKey) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSensitiveDoguConfigEntryRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockSensitiveDoguConfigEntryRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SensitiveDoguConfigKey
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *MockSensitiveDoguConfigEntryRepository_Get_Call {
	return &MockSensitiveDoguConfigEntryRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_Get_Call) Run(run func(_a0 context.Context, _a1 common.SensitiveDoguConfigKey)) *MockSensitiveDoguConfigEntryRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_Get_Call) Return(_a0 *ecosystem.SensitiveDoguConfigEntry, _a1 error) *MockSensitiveDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_Get_Call) RunAndReturn(run func(context.Context, common.SensitiveDoguConfigKey) (*ecosystem.SensitiveDoguConfigEntry, error)) *MockSensitiveDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByKey provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) GetAllByKey(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey) (map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAllByKey")
	}

	var r0 map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.SensitiveDoguConfigKey) (map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []common.SensitiveDoguConfigKey) map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []common.SensitiveDoguConfigKey) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByKey'
type MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call struct {
	*mock.Call
}

// GetAllByKey is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.SensitiveDoguConfigKey
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) GetAllByKey(_a0 interface{}, _a1 interface{}) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	return &MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", _a0, _a1)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) Run(run func(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey)) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) Return(_a0 map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, _a1 error) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []common.SensitiveDoguConfigKey) (map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, error)) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) Save(_a0 context.Context, _a1 *ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSensitiveDoguConfigEntryRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockSensitiveDoguConfigEntryRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ecosystem.SensitiveDoguConfigEntry
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *MockSensitiveDoguConfigEntryRepository_Save_Call {
	return &MockSensitiveDoguConfigEntryRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *ecosystem.SensitiveDoguConfigEntry)) *MockSensitiveDoguConfigEntryRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_Save_Call) Return(_a0 error) *MockSensitiveDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_Save_Call) RunAndReturn(run func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error) *MockSensitiveDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) SaveAll(_a0 context.Context, _a1 []*ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for SaveAll")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSensitiveDoguConfigEntryRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type MockSensitiveDoguConfigEntryRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []*ecosystem.SensitiveDoguConfigEntry
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) SaveAll(_a0 interface{}, _a1 interface{}) *MockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	return &MockSensitiveDoguConfigEntryRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", _a0, _a1)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveAll_Call) Run(run func(_a0 context.Context, _a1 []*ecosystem.SensitiveDoguConfigEntry)) *MockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveAll_Call) Return(_a0 error) *MockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []*ecosystem.SensitiveDoguConfigEntry) error) *MockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSensitiveDoguConfigEntryRepository creates a new instance of MockSensitiveDoguConfigEntryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSensitiveDoguConfigEntryRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSensitiveDoguConfigEntryRepository {
	mock := &MockSensitiveDoguConfigEntryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
