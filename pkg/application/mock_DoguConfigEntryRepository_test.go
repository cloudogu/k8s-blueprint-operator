// Code generated by mockery v2.20.0. DO NOT EDIT.

package application

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// MockDoguConfigEntryRepository is an autogenerated mock type for the DoguConfigEntryRepository type
type MockDoguConfigEntryRepository struct {
	mock.Mock
}

type MockDoguConfigEntryRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDoguConfigEntryRepository) EXPECT() *MockDoguConfigEntryRepository_Expecter {
	return &MockDoguConfigEntryRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *MockDoguConfigEntryRepository) Delete(_a0 context.Context, _a1 common.DoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.DoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguConfigEntryRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockDoguConfigEntryRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.DoguConfigKey
func (_e *MockDoguConfigEntryRepository_Expecter) Delete(_a0 interface{}, _a1 interface{}) *MockDoguConfigEntryRepository_Delete_Call {
	return &MockDoguConfigEntryRepository_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *MockDoguConfigEntryRepository_Delete_Call) Run(run func(_a0 context.Context, _a1 common.DoguConfigKey)) *MockDoguConfigEntryRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.DoguConfigKey))
	})
	return _c
}

func (_c *MockDoguConfigEntryRepository_Delete_Call) Return(_a0 error) *MockDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguConfigEntryRepository_Delete_Call) RunAndReturn(run func(context.Context, common.DoguConfigKey) error) *MockDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *MockDoguConfigEntryRepository) Get(_a0 context.Context, _a1 common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error) {
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

// MockDoguConfigEntryRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockDoguConfigEntryRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.DoguConfigKey
func (_e *MockDoguConfigEntryRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *MockDoguConfigEntryRepository_Get_Call {
	return &MockDoguConfigEntryRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *MockDoguConfigEntryRepository_Get_Call) Run(run func(_a0 context.Context, _a1 common.DoguConfigKey)) *MockDoguConfigEntryRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.DoguConfigKey))
	})
	return _c
}

func (_c *MockDoguConfigEntryRepository_Get_Call) Return(_a0 *ecosystem.DoguConfigEntry, _a1 error) *MockDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigEntryRepository_Get_Call) RunAndReturn(run func(context.Context, common.DoguConfigKey) (*ecosystem.DoguConfigEntry, error)) *MockDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByKey provides a mock function with given fields: _a0, _a1
func (_m *MockDoguConfigEntryRepository) GetAllByKey(_a0 context.Context, _a1 []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error) {
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

// MockDoguConfigEntryRepository_GetAllByKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByKey'
type MockDoguConfigEntryRepository_GetAllByKey_Call struct {
	*mock.Call
}

// GetAllByKey is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.DoguConfigKey
func (_e *MockDoguConfigEntryRepository_Expecter) GetAllByKey(_a0 interface{}, _a1 interface{}) *MockDoguConfigEntryRepository_GetAllByKey_Call {
	return &MockDoguConfigEntryRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", _a0, _a1)}
}

func (_c *MockDoguConfigEntryRepository_GetAllByKey_Call) Run(run func(_a0 context.Context, _a1 []common.DoguConfigKey)) *MockDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.DoguConfigKey))
	})
	return _c
}

func (_c *MockDoguConfigEntryRepository_GetAllByKey_Call) Return(_a0 map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, _a1 error) *MockDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguConfigEntryRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []common.DoguConfigKey) (map[common.DoguConfigKey]*ecosystem.DoguConfigEntry, error)) *MockDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *MockDoguConfigEntryRepository) Save(_a0 context.Context, _a1 *ecosystem.DoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.DoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguConfigEntryRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockDoguConfigEntryRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ecosystem.DoguConfigEntry
func (_e *MockDoguConfigEntryRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *MockDoguConfigEntryRepository_Save_Call {
	return &MockDoguConfigEntryRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *MockDoguConfigEntryRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *ecosystem.DoguConfigEntry)) *MockDoguConfigEntryRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.DoguConfigEntry))
	})
	return _c
}

func (_c *MockDoguConfigEntryRepository_Save_Call) Return(_a0 error) *MockDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguConfigEntryRepository_Save_Call) RunAndReturn(run func(context.Context, *ecosystem.DoguConfigEntry) error) *MockDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: _a0, _a1
func (_m *MockDoguConfigEntryRepository) SaveAll(_a0 context.Context, _a1 []*ecosystem.DoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.DoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguConfigEntryRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type MockDoguConfigEntryRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []*ecosystem.DoguConfigEntry
func (_e *MockDoguConfigEntryRepository_Expecter) SaveAll(_a0 interface{}, _a1 interface{}) *MockDoguConfigEntryRepository_SaveAll_Call {
	return &MockDoguConfigEntryRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", _a0, _a1)}
}

func (_c *MockDoguConfigEntryRepository_SaveAll_Call) Run(run func(_a0 context.Context, _a1 []*ecosystem.DoguConfigEntry)) *MockDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*ecosystem.DoguConfigEntry))
	})
	return _c
}

func (_c *MockDoguConfigEntryRepository_SaveAll_Call) Return(_a0 error) *MockDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguConfigEntryRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []*ecosystem.DoguConfigEntry) error) *MockDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDoguConfigEntryRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDoguConfigEntryRepository creates a new instance of MockDoguConfigEntryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDoguConfigEntryRepository(t mockConstructorTestingTNewMockDoguConfigEntryRepository) *MockDoguConfigEntryRepository {
	mock := &MockDoguConfigEntryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}