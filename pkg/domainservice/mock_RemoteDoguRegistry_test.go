// Code generated by mockery v2.42.1. DO NOT EDIT.

package domainservice

import (
	core "github.com/cloudogu/cesapp-lib/core"
	common "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"

	mock "github.com/stretchr/testify/mock"
)

// MockRemoteDoguRegistry is an autogenerated mock type for the RemoteDoguRegistry type
type MockRemoteDoguRegistry struct {
	mock.Mock
}

type MockRemoteDoguRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRemoteDoguRegistry) EXPECT() *MockRemoteDoguRegistry_Expecter {
	return &MockRemoteDoguRegistry_Expecter{mock: &_m.Mock}
}

// GetDogu provides a mock function with given fields: doguName, version
func (_m *MockRemoteDoguRegistry) GetDogu(doguName common.QualifiedDoguName, version string) (*core.Dogu, error) {
	ret := _m.Called(doguName, version)

	if len(ret) == 0 {
		panic("no return value specified for GetDogu")
	}

	var r0 *core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(common.QualifiedDoguName, string) (*core.Dogu, error)); ok {
		return rf(doguName, version)
	}
	if rf, ok := ret.Get(0).(func(common.QualifiedDoguName, string) *core.Dogu); ok {
		r0 = rf(doguName, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(common.QualifiedDoguName, string) error); ok {
		r1 = rf(doguName, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRemoteDoguRegistry_GetDogu_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDogu'
type MockRemoteDoguRegistry_GetDogu_Call struct {
	*mock.Call
}

// GetDogu is a helper method to define mock.On call
//   - doguName common.QualifiedDoguName
//   - version string
func (_e *MockRemoteDoguRegistry_Expecter) GetDogu(doguName interface{}, version interface{}) *MockRemoteDoguRegistry_GetDogu_Call {
	return &MockRemoteDoguRegistry_GetDogu_Call{Call: _e.mock.On("GetDogu", doguName, version)}
}

func (_c *MockRemoteDoguRegistry_GetDogu_Call) Run(run func(doguName common.QualifiedDoguName, version string)) *MockRemoteDoguRegistry_GetDogu_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.QualifiedDoguName), args[1].(string))
	})
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogu_Call) Return(_a0 *core.Dogu, _a1 error) *MockRemoteDoguRegistry_GetDogu_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogu_Call) RunAndReturn(run func(common.QualifiedDoguName, string) (*core.Dogu, error)) *MockRemoteDoguRegistry_GetDogu_Call {
	_c.Call.Return(run)
	return _c
}

// GetDogus provides a mock function with given fields: dogusToLoad
func (_m *MockRemoteDoguRegistry) GetDogus(dogusToLoad []DoguToLoad) (map[common.QualifiedDoguName]*core.Dogu, error) {
	ret := _m.Called(dogusToLoad)

	if len(ret) == 0 {
		panic("no return value specified for GetDogus")
	}

	var r0 map[common.QualifiedDoguName]*core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func([]DoguToLoad) (map[common.QualifiedDoguName]*core.Dogu, error)); ok {
		return rf(dogusToLoad)
	}
	if rf, ok := ret.Get(0).(func([]DoguToLoad) map[common.QualifiedDoguName]*core.Dogu); ok {
		r0 = rf(dogusToLoad)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.QualifiedDoguName]*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func([]DoguToLoad) error); ok {
		r1 = rf(dogusToLoad)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRemoteDoguRegistry_GetDogus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDogus'
type MockRemoteDoguRegistry_GetDogus_Call struct {
	*mock.Call
}

// GetDogus is a helper method to define mock.On call
//   - dogusToLoad []DoguToLoad
func (_e *MockRemoteDoguRegistry_Expecter) GetDogus(dogusToLoad interface{}) *MockRemoteDoguRegistry_GetDogus_Call {
	return &MockRemoteDoguRegistry_GetDogus_Call{Call: _e.mock.On("GetDogus", dogusToLoad)}
}

func (_c *MockRemoteDoguRegistry_GetDogus_Call) Run(run func(dogusToLoad []DoguToLoad)) *MockRemoteDoguRegistry_GetDogus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]DoguToLoad))
	})
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogus_Call) Return(_a0 map[common.QualifiedDoguName]*core.Dogu, _a1 error) *MockRemoteDoguRegistry_GetDogus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogus_Call) RunAndReturn(run func([]DoguToLoad) (map[common.QualifiedDoguName]*core.Dogu, error)) *MockRemoteDoguRegistry_GetDogus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRemoteDoguRegistry creates a new instance of MockRemoteDoguRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRemoteDoguRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRemoteDoguRegistry {
	mock := &MockRemoteDoguRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
