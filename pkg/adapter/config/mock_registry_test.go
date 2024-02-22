// Code generated by mockery v2.20.0. DO NOT EDIT.

package config

import (
	registry "github.com/cloudogu/cesapp-lib/registry"
	mock "github.com/stretchr/testify/mock"
)

// mockRegistry is an autogenerated mock type for the registry type
type mockRegistry struct {
	mock.Mock
}

type mockRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *mockRegistry) EXPECT() *mockRegistry_Expecter {
	return &mockRegistry_Expecter{mock: &_m.Mock}
}

// BlueprintRegistry provides a mock function with given fields:
func (_m *mockRegistry) BlueprintRegistry() registry.ConfigurationContext {
	ret := _m.Called()

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func() registry.ConfigurationContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockRegistry_BlueprintRegistry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlueprintRegistry'
type mockRegistry_BlueprintRegistry_Call struct {
	*mock.Call
}

// BlueprintRegistry is a helper method to define mock.On call
func (_e *mockRegistry_Expecter) BlueprintRegistry() *mockRegistry_BlueprintRegistry_Call {
	return &mockRegistry_BlueprintRegistry_Call{Call: _e.mock.On("BlueprintRegistry")}
}

func (_c *mockRegistry_BlueprintRegistry_Call) Run(run func()) *mockRegistry_BlueprintRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRegistry_BlueprintRegistry_Call) Return(_a0 registry.ConfigurationContext) *mockRegistry_BlueprintRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_BlueprintRegistry_Call) RunAndReturn(run func() registry.ConfigurationContext) *mockRegistry_BlueprintRegistry_Call {
	_c.Call.Return(run)
	return _c
}

// DoguConfig provides a mock function with given fields: dogu
func (_m *mockRegistry) DoguConfig(dogu string) registry.ConfigurationContext {
	ret := _m.Called(dogu)

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func(string) registry.ConfigurationContext); ok {
		r0 = rf(dogu)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockRegistry_DoguConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DoguConfig'
type mockRegistry_DoguConfig_Call struct {
	*mock.Call
}

// DoguConfig is a helper method to define mock.On call
//   - dogu string
func (_e *mockRegistry_Expecter) DoguConfig(dogu interface{}) *mockRegistry_DoguConfig_Call {
	return &mockRegistry_DoguConfig_Call{Call: _e.mock.On("DoguConfig", dogu)}
}

func (_c *mockRegistry_DoguConfig_Call) Run(run func(dogu string)) *mockRegistry_DoguConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockRegistry_DoguConfig_Call) Return(_a0 registry.ConfigurationContext) *mockRegistry_DoguConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_DoguConfig_Call) RunAndReturn(run func(string) registry.ConfigurationContext) *mockRegistry_DoguConfig_Call {
	_c.Call.Return(run)
	return _c
}

// DoguRegistry provides a mock function with given fields:
func (_m *mockRegistry) DoguRegistry() registry.DoguRegistry {
	ret := _m.Called()

	var r0 registry.DoguRegistry
	if rf, ok := ret.Get(0).(func() registry.DoguRegistry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.DoguRegistry)
		}
	}

	return r0
}

// mockRegistry_DoguRegistry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DoguRegistry'
type mockRegistry_DoguRegistry_Call struct {
	*mock.Call
}

// DoguRegistry is a helper method to define mock.On call
func (_e *mockRegistry_Expecter) DoguRegistry() *mockRegistry_DoguRegistry_Call {
	return &mockRegistry_DoguRegistry_Call{Call: _e.mock.On("DoguRegistry")}
}

func (_c *mockRegistry_DoguRegistry_Call) Run(run func()) *mockRegistry_DoguRegistry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRegistry_DoguRegistry_Call) Return(_a0 registry.DoguRegistry) *mockRegistry_DoguRegistry_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_DoguRegistry_Call) RunAndReturn(run func() registry.DoguRegistry) *mockRegistry_DoguRegistry_Call {
	_c.Call.Return(run)
	return _c
}

// GetNode provides a mock function with given fields:
func (_m *mockRegistry) GetNode() (registry.Node, error) {
	ret := _m.Called()

	var r0 registry.Node
	var r1 error
	if rf, ok := ret.Get(0).(func() (registry.Node, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() registry.Node); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(registry.Node)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockRegistry_GetNode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetNode'
type mockRegistry_GetNode_Call struct {
	*mock.Call
}

// GetNode is a helper method to define mock.On call
func (_e *mockRegistry_Expecter) GetNode() *mockRegistry_GetNode_Call {
	return &mockRegistry_GetNode_Call{Call: _e.mock.On("GetNode")}
}

func (_c *mockRegistry_GetNode_Call) Run(run func()) *mockRegistry_GetNode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRegistry_GetNode_Call) Return(_a0 registry.Node, _a1 error) *mockRegistry_GetNode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockRegistry_GetNode_Call) RunAndReturn(run func() (registry.Node, error)) *mockRegistry_GetNode_Call {
	_c.Call.Return(run)
	return _c
}

// GlobalConfig provides a mock function with given fields:
func (_m *mockRegistry) GlobalConfig() registry.ConfigurationContext {
	ret := _m.Called()

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func() registry.ConfigurationContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockRegistry_GlobalConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GlobalConfig'
type mockRegistry_GlobalConfig_Call struct {
	*mock.Call
}

// GlobalConfig is a helper method to define mock.On call
func (_e *mockRegistry_Expecter) GlobalConfig() *mockRegistry_GlobalConfig_Call {
	return &mockRegistry_GlobalConfig_Call{Call: _e.mock.On("GlobalConfig")}
}

func (_c *mockRegistry_GlobalConfig_Call) Run(run func()) *mockRegistry_GlobalConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRegistry_GlobalConfig_Call) Return(_a0 registry.ConfigurationContext) *mockRegistry_GlobalConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_GlobalConfig_Call) RunAndReturn(run func() registry.ConfigurationContext) *mockRegistry_GlobalConfig_Call {
	_c.Call.Return(run)
	return _c
}

// HostConfig provides a mock function with given fields: hostService
func (_m *mockRegistry) HostConfig(hostService string) registry.ConfigurationContext {
	ret := _m.Called(hostService)

	var r0 registry.ConfigurationContext
	if rf, ok := ret.Get(0).(func(string) registry.ConfigurationContext); ok {
		r0 = rf(hostService)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.ConfigurationContext)
		}
	}

	return r0
}

// mockRegistry_HostConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HostConfig'
type mockRegistry_HostConfig_Call struct {
	*mock.Call
}

// HostConfig is a helper method to define mock.On call
//   - hostService string
func (_e *mockRegistry_Expecter) HostConfig(hostService interface{}) *mockRegistry_HostConfig_Call {
	return &mockRegistry_HostConfig_Call{Call: _e.mock.On("HostConfig", hostService)}
}

func (_c *mockRegistry_HostConfig_Call) Run(run func(hostService string)) *mockRegistry_HostConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockRegistry_HostConfig_Call) Return(_a0 registry.ConfigurationContext) *mockRegistry_HostConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_HostConfig_Call) RunAndReturn(run func(string) registry.ConfigurationContext) *mockRegistry_HostConfig_Call {
	_c.Call.Return(run)
	return _c
}

// RootConfig provides a mock function with given fields:
func (_m *mockRegistry) RootConfig() registry.WatchConfigurationContext {
	ret := _m.Called()

	var r0 registry.WatchConfigurationContext
	if rf, ok := ret.Get(0).(func() registry.WatchConfigurationContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.WatchConfigurationContext)
		}
	}

	return r0
}

// mockRegistry_RootConfig_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RootConfig'
type mockRegistry_RootConfig_Call struct {
	*mock.Call
}

// RootConfig is a helper method to define mock.On call
func (_e *mockRegistry_Expecter) RootConfig() *mockRegistry_RootConfig_Call {
	return &mockRegistry_RootConfig_Call{Call: _e.mock.On("RootConfig")}
}

func (_c *mockRegistry_RootConfig_Call) Run(run func()) *mockRegistry_RootConfig_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRegistry_RootConfig_Call) Return(_a0 registry.WatchConfigurationContext) *mockRegistry_RootConfig_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_RootConfig_Call) RunAndReturn(run func() registry.WatchConfigurationContext) *mockRegistry_RootConfig_Call {
	_c.Call.Return(run)
	return _c
}

// State provides a mock function with given fields: dogu
func (_m *mockRegistry) State(dogu string) registry.State {
	ret := _m.Called(dogu)

	var r0 registry.State
	if rf, ok := ret.Get(0).(func(string) registry.State); ok {
		r0 = rf(dogu)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(registry.State)
		}
	}

	return r0
}

// mockRegistry_State_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'State'
type mockRegistry_State_Call struct {
	*mock.Call
}

// State is a helper method to define mock.On call
//   - dogu string
func (_e *mockRegistry_Expecter) State(dogu interface{}) *mockRegistry_State_Call {
	return &mockRegistry_State_Call{Call: _e.mock.On("State", dogu)}
}

func (_c *mockRegistry_State_Call) Run(run func(dogu string)) *mockRegistry_State_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockRegistry_State_Call) Return(_a0 registry.State) *mockRegistry_State_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRegistry_State_Call) RunAndReturn(run func(string) registry.State) *mockRegistry_State_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockRegistry interface {
	mock.TestingT
	Cleanup(func())
}

// newMockRegistry creates a new instance of mockRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockRegistry(t mockConstructorTestingTnewMockRegistry) *mockRegistry {
	mock := &mockRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}