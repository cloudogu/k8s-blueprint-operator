// Code generated by mockery v2.20.0. DO NOT EDIT.

package application

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// mockSensitiveDoguConfigEntryRepository is an autogenerated mock type for the sensitiveDoguConfigEntryRepository type
type mockSensitiveDoguConfigEntryRepository struct {
	mock.Mock
}

type mockSensitiveDoguConfigEntryRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockSensitiveDoguConfigEntryRepository) EXPECT() *mockSensitiveDoguConfigEntryRepository_Expecter {
	return &mockSensitiveDoguConfigEntryRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *mockSensitiveDoguConfigEntryRepository) Delete(_a0 context.Context, _a1 common.SensitiveDoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SensitiveDoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSensitiveDoguConfigEntryRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockSensitiveDoguConfigEntryRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SensitiveDoguConfigKey
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) Delete(_a0 interface{}, _a1 interface{}) *mockSensitiveDoguConfigEntryRepository_Delete_Call {
	return &mockSensitiveDoguConfigEntryRepository_Delete_Call{Call: _e.mock.On("Delete", _a0, _a1)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_Delete_Call) Run(run func(_a0 context.Context, _a1 common.SensitiveDoguConfigKey)) *mockSensitiveDoguConfigEntryRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_Delete_Call) Return(_a0 error) *mockSensitiveDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_Delete_Call) RunAndReturn(run func(context.Context, common.SensitiveDoguConfigKey) error) *mockSensitiveDoguConfigEntryRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAllByKeys provides a mock function with given fields: _a0, _a1
func (_m *mockSensitiveDoguConfigEntryRepository) DeleteAllByKeys(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.SensitiveDoguConfigKey) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAllByKeys'
type mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call struct {
	*mock.Call
}

// DeleteAllByKeys is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.SensitiveDoguConfigKey
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) DeleteAllByKeys(_a0 interface{}, _a1 interface{}) *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	return &mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call{Call: _e.mock.On("DeleteAllByKeys", _a0, _a1)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call) Run(run func(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey)) *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call) Return(_a0 error) *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call) RunAndReturn(run func(context.Context, []common.SensitiveDoguConfigKey) error) *mockSensitiveDoguConfigEntryRepository_DeleteAllByKeys_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *mockSensitiveDoguConfigEntryRepository) Get(_a0 context.Context, _a1 common.SensitiveDoguConfigKey) (*ecosystem.SensitiveDoguConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

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

// mockSensitiveDoguConfigEntryRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockSensitiveDoguConfigEntryRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SensitiveDoguConfigKey
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) Get(_a0 interface{}, _a1 interface{}) *mockSensitiveDoguConfigEntryRepository_Get_Call {
	return &mockSensitiveDoguConfigEntryRepository_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_Get_Call) Run(run func(_a0 context.Context, _a1 common.SensitiveDoguConfigKey)) *mockSensitiveDoguConfigEntryRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_Get_Call) Return(_a0 *ecosystem.SensitiveDoguConfigEntry, _a1 error) *mockSensitiveDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_Get_Call) RunAndReturn(run func(context.Context, common.SensitiveDoguConfigKey) (*ecosystem.SensitiveDoguConfigEntry, error)) *mockSensitiveDoguConfigEntryRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByKey provides a mock function with given fields: _a0, _a1
func (_m *mockSensitiveDoguConfigEntryRepository) GetAllByKey(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey) (map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, error) {
	ret := _m.Called(_a0, _a1)

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

// mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByKey'
type mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call struct {
	*mock.Call
}

// GetAllByKey is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []common.SensitiveDoguConfigKey
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) GetAllByKey(_a0 interface{}, _a1 interface{}) *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	return &mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", _a0, _a1)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) Run(run func(_a0 context.Context, _a1 []common.SensitiveDoguConfigKey)) *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) Return(_a0 map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, _a1 error) *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []common.SensitiveDoguConfigKey) (map[common.SensitiveDoguConfigKey]*ecosystem.SensitiveDoguConfigEntry, error)) *mockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *mockSensitiveDoguConfigEntryRepository) Save(_a0 context.Context, _a1 *ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSensitiveDoguConfigEntryRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type mockSensitiveDoguConfigEntryRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ecosystem.SensitiveDoguConfigEntry
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *mockSensitiveDoguConfigEntryRepository_Save_Call {
	return &mockSensitiveDoguConfigEntryRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *ecosystem.SensitiveDoguConfigEntry)) *mockSensitiveDoguConfigEntryRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_Save_Call) Return(_a0 error) *mockSensitiveDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_Save_Call) RunAndReturn(run func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error) *mockSensitiveDoguConfigEntryRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: _a0, _a1
func (_m *mockSensitiveDoguConfigEntryRepository) SaveAll(_a0 context.Context, _a1 []*ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSensitiveDoguConfigEntryRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type mockSensitiveDoguConfigEntryRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 []*ecosystem.SensitiveDoguConfigEntry
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) SaveAll(_a0 interface{}, _a1 interface{}) *mockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	return &mockSensitiveDoguConfigEntryRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", _a0, _a1)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveAll_Call) Run(run func(_a0 context.Context, _a1 []*ecosystem.SensitiveDoguConfigEntry)) *mockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveAll_Call) Return(_a0 error) *mockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []*ecosystem.SensitiveDoguConfigEntry) error) *mockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAllForNotInstalledDogu provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockSensitiveDoguConfigEntryRepository) SaveAllForNotInstalledDogu(_a0 context.Context, _a1 common.SimpleDoguName, _a2 []*ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleDoguName, []*ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAllForNotInstalledDogu'
type mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call struct {
	*mock.Call
}

// SaveAllForNotInstalledDogu is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SimpleDoguName
//   - _a2 []*ecosystem.SensitiveDoguConfigEntry
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) SaveAllForNotInstalledDogu(_a0 interface{}, _a1 interface{}, _a2 interface{}) *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call {
	return &mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call{Call: _e.mock.On("SaveAllForNotInstalledDogu", _a0, _a1, _a2)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call) Run(run func(_a0 context.Context, _a1 common.SimpleDoguName, _a2 []*ecosystem.SensitiveDoguConfigEntry)) *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleDoguName), args[2].([]*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call) Return(_a0 error) *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call) RunAndReturn(run func(context.Context, common.SimpleDoguName, []*ecosystem.SensitiveDoguConfigEntry) error) *mockSensitiveDoguConfigEntryRepository_SaveAllForNotInstalledDogu_Call {
	_c.Call.Return(run)
	return _c
}

// SaveForNotInstalledDogu provides a mock function with given fields: ctx, entry
func (_m *mockSensitiveDoguConfigEntryRepository) SaveForNotInstalledDogu(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(ctx, entry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(ctx, entry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveForNotInstalledDogu'
type mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call struct {
	*mock.Call
}

// SaveForNotInstalledDogu is a helper method to define mock.On call
//   - ctx context.Context
//   - entry *ecosystem.SensitiveDoguConfigEntry
func (_e *mockSensitiveDoguConfigEntryRepository_Expecter) SaveForNotInstalledDogu(ctx interface{}, entry interface{}) *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	return &mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call{Call: _e.mock.On("SaveForNotInstalledDogu", ctx, entry)}
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call) Run(run func(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry)) *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call) Return(_a0 error) *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call) RunAndReturn(run func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error) *mockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockSensitiveDoguConfigEntryRepository interface {
	mock.TestingT
	Cleanup(func())
}

// newMockSensitiveDoguConfigEntryRepository creates a new instance of mockSensitiveDoguConfigEntryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockSensitiveDoguConfigEntryRepository(t mockConstructorTestingTnewMockSensitiveDoguConfigEntryRepository) *mockSensitiveDoguConfigEntryRepository {
	mock := &mockSensitiveDoguConfigEntryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
