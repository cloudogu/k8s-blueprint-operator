// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// MockDoguSensibleConfigKeyRepository is an autogenerated mock type for the DoguSensibleConfigKeyRepository type
type MockDoguSensibleConfigKeyRepository struct {
	mock.Mock
}

type MockDoguSensibleConfigKeyRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDoguSensibleConfigKeyRepository) EXPECT() *MockDoguSensibleConfigKeyRepository_Expecter {
	return &MockDoguSensibleConfigKeyRepository_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, key
func (_m *MockDoguSensibleConfigKeyRepository) Delete(ctx context.Context, key ecosystem.DoguConfigKey) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ecosystem.DoguConfigKey) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguSensibleConfigKeyRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockDoguSensibleConfigKeyRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - key ecosystem.DoguConfigKey
func (_e *MockDoguSensibleConfigKeyRepository_Expecter) Delete(ctx interface{}, key interface{}) *MockDoguSensibleConfigKeyRepository_Delete_Call {
	return &MockDoguSensibleConfigKeyRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, key)}
}

func (_c *MockDoguSensibleConfigKeyRepository_Delete_Call) Run(run func(ctx context.Context, key ecosystem.DoguConfigKey)) *MockDoguSensibleConfigKeyRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ecosystem.DoguConfigKey))
	})
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_Delete_Call) Return(_a0 error) *MockDoguSensibleConfigKeyRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_Delete_Call) RunAndReturn(run func(context.Context, ecosystem.DoguConfigKey) error) *MockDoguSensibleConfigKeyRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllByKey provides a mock function with given fields: ctx, keys
func (_m *MockDoguSensibleConfigKeyRepository) GetAllByKey(ctx context.Context, keys []ecosystem.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error) {
	ret := _m.Called(ctx, keys)

	var r0 map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []ecosystem.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error)); ok {
		return rf(ctx, keys)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []ecosystem.DoguConfigKey) map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry); ok {
		r0 = rf(ctx, keys)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []ecosystem.DoguConfigKey) error); ok {
		r1 = rf(ctx, keys)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguSensibleConfigKeyRepository_GetAllByKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllByKey'
type MockDoguSensibleConfigKeyRepository_GetAllByKey_Call struct {
	*mock.Call
}

// GetAllByKey is a helper method to define mock.On call
//   - ctx context.Context
//   - keys []ecosystem.DoguConfigKey
func (_e *MockDoguSensibleConfigKeyRepository_Expecter) GetAllByKey(ctx interface{}, keys interface{}) *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call {
	return &MockDoguSensibleConfigKeyRepository_GetAllByKey_Call{Call: _e.mock.On("GetAllByKey", ctx, keys)}
}

func (_c *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call) Run(run func(ctx context.Context, keys []ecosystem.DoguConfigKey)) *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]ecosystem.DoguConfigKey))
	})
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call) Return(_a0 map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, _a1 error) *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call) RunAndReturn(run func(context.Context, []ecosystem.DoguConfigKey) (map[common.SimpleDoguName][]*ecosystem.DoguConfigEntry, error)) *MockDoguSensibleConfigKeyRepository_GetAllByKey_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: _a0, _a1
func (_m *MockDoguSensibleConfigKeyRepository) Save(_a0 context.Context, _a1 *ecosystem.DoguConfigEntry) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.DoguConfigEntry) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguSensibleConfigKeyRepository_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type MockDoguSensibleConfigKeyRepository_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *ecosystem.DoguConfigEntry
func (_e *MockDoguSensibleConfigKeyRepository_Expecter) Save(_a0 interface{}, _a1 interface{}) *MockDoguSensibleConfigKeyRepository_Save_Call {
	return &MockDoguSensibleConfigKeyRepository_Save_Call{Call: _e.mock.On("Save", _a0, _a1)}
}

func (_c *MockDoguSensibleConfigKeyRepository_Save_Call) Run(run func(_a0 context.Context, _a1 *ecosystem.DoguConfigEntry)) *MockDoguSensibleConfigKeyRepository_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.DoguConfigEntry))
	})
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_Save_Call) Return(_a0 error) *MockDoguSensibleConfigKeyRepository_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_Save_Call) RunAndReturn(run func(context.Context, *ecosystem.DoguConfigEntry) error) *MockDoguSensibleConfigKeyRepository_Save_Call {
	_c.Call.Return(run)
	return _c
}

// SaveAll provides a mock function with given fields: ctx, keys
func (_m *MockDoguSensibleConfigKeyRepository) SaveAll(ctx context.Context, keys []*ecosystem.DoguConfigEntry) error {
	ret := _m.Called(ctx, keys)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*ecosystem.DoguConfigEntry) error); ok {
		r0 = rf(ctx, keys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguSensibleConfigKeyRepository_SaveAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SaveAll'
type MockDoguSensibleConfigKeyRepository_SaveAll_Call struct {
	*mock.Call
}

// SaveAll is a helper method to define mock.On call
//   - ctx context.Context
//   - keys []*ecosystem.DoguConfigEntry
func (_e *MockDoguSensibleConfigKeyRepository_Expecter) SaveAll(ctx interface{}, keys interface{}) *MockDoguSensibleConfigKeyRepository_SaveAll_Call {
	return &MockDoguSensibleConfigKeyRepository_SaveAll_Call{Call: _e.mock.On("SaveAll", ctx, keys)}
}

func (_c *MockDoguSensibleConfigKeyRepository_SaveAll_Call) Run(run func(ctx context.Context, keys []*ecosystem.DoguConfigEntry)) *MockDoguSensibleConfigKeyRepository_SaveAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*ecosystem.DoguConfigEntry))
	})
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_SaveAll_Call) Return(_a0 error) *MockDoguSensibleConfigKeyRepository_SaveAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguSensibleConfigKeyRepository_SaveAll_Call) RunAndReturn(run func(context.Context, []*ecosystem.DoguConfigEntry) error) *MockDoguSensibleConfigKeyRepository_SaveAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDoguSensibleConfigKeyRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDoguSensibleConfigKeyRepository creates a new instance of MockDoguSensibleConfigKeyRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDoguSensibleConfigKeyRepository(t mockConstructorTestingTNewMockDoguSensibleConfigKeyRepository) *MockDoguSensibleConfigKeyRepository {
	mock := &MockDoguSensibleConfigKeyRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
