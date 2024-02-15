// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	mock "github.com/stretchr/testify/mock"
)

// MockConfigEncryptionAdapter is an autogenerated mock type for the ConfigEncryptionAdapter type
type MockConfigEncryptionAdapter struct {
	mock.Mock
}

type MockConfigEncryptionAdapter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockConfigEncryptionAdapter) EXPECT() *MockConfigEncryptionAdapter_Expecter {
	return &MockConfigEncryptionAdapter_Expecter{mock: &_m.Mock}
}

// Encrypt provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockConfigEncryptionAdapter) Encrypt(_a0 context.Context, _a1 common.SimpleDoguName, _a2 common.SensitiveDoguConfigValue) (common.EncryptedDoguConfigValue, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 common.EncryptedDoguConfigValue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleDoguName, common.SensitiveDoguConfigValue) (common.EncryptedDoguConfigValue, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleDoguName, common.SensitiveDoguConfigValue) common.EncryptedDoguConfigValue); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Get(0).(common.EncryptedDoguConfigValue)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.SimpleDoguName, common.SensitiveDoguConfigValue) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConfigEncryptionAdapter_Encrypt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Encrypt'
type MockConfigEncryptionAdapter_Encrypt_Call struct {
	*mock.Call
}

// Encrypt is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SimpleDoguName
//   - _a2 common.SensitiveDoguConfigValue
func (_e *MockConfigEncryptionAdapter_Expecter) Encrypt(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockConfigEncryptionAdapter_Encrypt_Call {
	return &MockConfigEncryptionAdapter_Encrypt_Call{Call: _e.mock.On("Encrypt", _a0, _a1, _a2)}
}

func (_c *MockConfigEncryptionAdapter_Encrypt_Call) Run(run func(_a0 context.Context, _a1 common.SimpleDoguName, _a2 common.SensitiveDoguConfigValue)) *MockConfigEncryptionAdapter_Encrypt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleDoguName), args[2].(common.SensitiveDoguConfigValue))
	})
	return _c
}

func (_c *MockConfigEncryptionAdapter_Encrypt_Call) Return(_a0 common.EncryptedDoguConfigValue, _a1 error) *MockConfigEncryptionAdapter_Encrypt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConfigEncryptionAdapter_Encrypt_Call) RunAndReturn(run func(context.Context, common.SimpleDoguName, common.SensitiveDoguConfigValue) (common.EncryptedDoguConfigValue, error)) *MockConfigEncryptionAdapter_Encrypt_Call {
	_c.Call.Return(run)
	return _c
}

// EncryptAll provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockConfigEncryptionAdapter) EncryptAll(_a0 context.Context, _a1 common.SimpleDoguName, _a2 []common.SensitiveDoguConfigValue) (map[common.SimpleDoguName]common.EncryptedDoguConfigValue, error) {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 map[common.SimpleDoguName]common.EncryptedDoguConfigValue
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleDoguName, []common.SensitiveDoguConfigValue) (map[common.SimpleDoguName]common.EncryptedDoguConfigValue, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleDoguName, []common.SensitiveDoguConfigValue) map[common.SimpleDoguName]common.EncryptedDoguConfigValue); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.SimpleDoguName]common.EncryptedDoguConfigValue)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.SimpleDoguName, []common.SensitiveDoguConfigValue) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockConfigEncryptionAdapter_EncryptAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EncryptAll'
type MockConfigEncryptionAdapter_EncryptAll_Call struct {
	*mock.Call
}

// EncryptAll is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 common.SimpleDoguName
//   - _a2 []common.SensitiveDoguConfigValue
func (_e *MockConfigEncryptionAdapter_Expecter) EncryptAll(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockConfigEncryptionAdapter_EncryptAll_Call {
	return &MockConfigEncryptionAdapter_EncryptAll_Call{Call: _e.mock.On("EncryptAll", _a0, _a1, _a2)}
}

func (_c *MockConfigEncryptionAdapter_EncryptAll_Call) Run(run func(_a0 context.Context, _a1 common.SimpleDoguName, _a2 []common.SensitiveDoguConfigValue)) *MockConfigEncryptionAdapter_EncryptAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleDoguName), args[2].([]common.SensitiveDoguConfigValue))
	})
	return _c
}

func (_c *MockConfigEncryptionAdapter_EncryptAll_Call) Return(_a0 map[common.SimpleDoguName]common.EncryptedDoguConfigValue, _a1 error) *MockConfigEncryptionAdapter_EncryptAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockConfigEncryptionAdapter_EncryptAll_Call) RunAndReturn(run func(context.Context, common.SimpleDoguName, []common.SensitiveDoguConfigValue) (map[common.SimpleDoguName]common.EncryptedDoguConfigValue, error)) *MockConfigEncryptionAdapter_EncryptAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockConfigEncryptionAdapter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockConfigEncryptionAdapter creates a new instance of MockConfigEncryptionAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockConfigEncryptionAdapter(t mockConstructorTestingTNewMockConfigEncryptionAdapter) *MockConfigEncryptionAdapter {
	mock := &MockConfigEncryptionAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
