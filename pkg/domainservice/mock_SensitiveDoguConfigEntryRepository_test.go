// Code generated by mockery v2.20.0. DO NOT EDIT.

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

// Delete provides a mock function with given fields: ctx, key
func (_m *MockSensitiveDoguConfigEntryRepository) Delete(ctx context.Context, key common.SensitiveDoguConfigKey) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SensitiveDoguConfigKey) error); ok {
		r0 = rf(ctx, key)
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
//   - ctx context.Context
//   - key common.SensitiveDoguConfigKey
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) Delete(ctx interface{}, key interface{}) *MockSensitiveDoguConfigEntryRepository_Delete_Call {
	return &MockSensitiveDoguConfigEntryRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, key)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_Delete_Call) Run(run func(ctx context.Context, key common.SensitiveDoguConfigKey)) *MockSensitiveDoguConfigEntryRepository_Delete_Call {
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

// GetAllByKey provides a mock function with given fields: ctx, keys
func (_m *MockSensitiveDoguConfigEntryRepository) GetAllByKey(ctx context.Context, keys []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error) {
	ret := _m.Called(ctx, keys)

	var r0 map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error)); ok {
		return rf(ctx, keys)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []common.SensitiveDoguConfigKey) map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry); ok {
		r0 = rf(ctx, keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []common.SensitiveDoguConfigKey) error); ok {
		r1 = rf(ctx, keys)
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
//   - ctx context.Context
//   - keys []common.SensitiveDoguConfigKey
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) GetAllByKey(ctx interface{}, keys interface{}) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	return &MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", ctx, keys)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) Run(run func(ctx context.Context, keys []common.SensitiveDoguConfigKey)) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.SensitiveDoguConfigKey))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) Return(_a0 map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, _a1 error) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []common.SensitiveDoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.SensitiveDoguConfigEntry, error)) *MockSensitiveDoguConfigEntryRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *MockSensitiveDoguConfigEntryRepository) Save(_a0 context.Context, _a1 *ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

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

// SaveAll provides a mock function with given fields: ctx, keys
func (_m *MockSensitiveDoguConfigEntryRepository) SaveAll(ctx context.Context, keys []*ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(ctx, keys)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(ctx, keys)
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
//   - ctx context.Context
//   - keys []*ecosystem.SensitiveDoguConfigEntry
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) SaveAll(ctx interface{}, keys interface{}) *MockSensitiveDoguConfigEntryRepository_SaveAll_Call {
	return &MockSensitiveDoguConfigEntryRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", ctx, keys)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveAll_Call) Run(run func(ctx context.Context, keys []*ecosystem.SensitiveDoguConfigEntry)) *MockSensitiveDoguConfigEntryRepository_SaveAll_Call {
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

// SaveForNotInstalledDogu provides a mock function with given fields: ctx, entry
func (_m *MockSensitiveDoguConfigEntryRepository) SaveForNotInstalledDogu(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry) error {
	ret := _m.Called(ctx, entry)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error); ok {
		r0 = rf(ctx, entry)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveForNotInstalledDogu'
type MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call struct {
	*mock.Call
}

// SaveForNotInstalledDogu is a helper method to define mock.On call
//   - ctx context.Context
//   - entry *ecosystem.SensitiveDoguConfigEntry
func (_e *MockSensitiveDoguConfigEntryRepository_Expecter) SaveForNotInstalledDogu(ctx interface{}, entry interface{}) *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	return &MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call{Call: _e.mock.On("SaveForNotInstalledDogu", ctx, entry)}
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call) Run(run func(ctx context.Context, entry *ecosystem.SensitiveDoguConfigEntry)) *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.SensitiveDoguConfigEntry))
	})
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call) Return(_a0 error) *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call) RunAndReturn(run func(context.Context, *ecosystem.SensitiveDoguConfigEntry) error) *MockSensitiveDoguConfigEntryRepository_SaveForNotInstalledDogu_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockSensitiveDoguConfigEntryRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSensitiveDoguConfigEntryRepository creates a new instance of MockSensitiveDoguConfigEntryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSensitiveDoguConfigEntryRepository(t mockConstructorTestingTNewMockSensitiveDoguConfigEntryRepository) *MockSensitiveDoguConfigEntryRepository {
	mock := &MockSensitiveDoguConfigEntryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
