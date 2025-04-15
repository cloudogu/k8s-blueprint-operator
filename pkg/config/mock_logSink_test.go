// Code generated by mockery v2.53.3. DO NOT EDIT.

package config

import (
	logr "github.com/go-logr/logr"
	mock "github.com/stretchr/testify/mock"
)

// mockLogSink is an autogenerated mock type for the logSink type
type mockLogSink struct {
	mock.Mock
}

type mockLogSink_Expecter struct {
	mock *mock.Mock
}

func (_m *mockLogSink) EXPECT() *mockLogSink_Expecter {
	return &mockLogSink_Expecter{mock: &_m.Mock}
}

// Enabled provides a mock function with given fields: level
func (_m *mockLogSink) Enabled(level int) bool {
	ret := _m.Called(level)

	if len(ret) == 0 {
		panic("no return value specified for Enabled")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(int) bool); ok {
		r0 = rf(level)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// mockLogSink_Enabled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Enabled'
type mockLogSink_Enabled_Call struct {
	*mock.Call
}

// Enabled is a helper method to define mock.On call
//   - level int
func (_e *mockLogSink_Expecter) Enabled(level interface{}) *mockLogSink_Enabled_Call {
	return &mockLogSink_Enabled_Call{Call: _e.mock.On("Enabled", level)}
}

func (_c *mockLogSink_Enabled_Call) Run(run func(level int)) *mockLogSink_Enabled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *mockLogSink_Enabled_Call) Return(_a0 bool) *mockLogSink_Enabled_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockLogSink_Enabled_Call) RunAndReturn(run func(int) bool) *mockLogSink_Enabled_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with given fields: err, msg, keysAndValues
func (_m *mockLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, err, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// mockLogSink_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type mockLogSink_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
//   - err error
//   - msg string
//   - keysAndValues ...interface{}
func (_e *mockLogSink_Expecter) Error(err interface{}, msg interface{}, keysAndValues ...interface{}) *mockLogSink_Error_Call {
	return &mockLogSink_Error_Call{Call: _e.mock.On("Error",
		append([]interface{}{err, msg}, keysAndValues...)...)}
}

func (_c *mockLogSink_Error_Call) Run(run func(err error, msg string, keysAndValues ...interface{})) *mockLogSink_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(error), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockLogSink_Error_Call) Return() *mockLogSink_Error_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogSink_Error_Call) RunAndReturn(run func(error, string, ...interface{})) *mockLogSink_Error_Call {
	_c.Run(run)
	return _c
}

// Info provides a mock function with given fields: level, msg, keysAndValues
func (_m *mockLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, level, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// mockLogSink_Info_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Info'
type mockLogSink_Info_Call struct {
	*mock.Call
}

// Info is a helper method to define mock.On call
//   - level int
//   - msg string
//   - keysAndValues ...interface{}
func (_e *mockLogSink_Expecter) Info(level interface{}, msg interface{}, keysAndValues ...interface{}) *mockLogSink_Info_Call {
	return &mockLogSink_Info_Call{Call: _e.mock.On("Info",
		append([]interface{}{level, msg}, keysAndValues...)...)}
}

func (_c *mockLogSink_Info_Call) Run(run func(level int, msg string, keysAndValues ...interface{})) *mockLogSink_Info_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(int), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *mockLogSink_Info_Call) Return() *mockLogSink_Info_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogSink_Info_Call) RunAndReturn(run func(int, string, ...interface{})) *mockLogSink_Info_Call {
	_c.Run(run)
	return _c
}

// Init provides a mock function with given fields: info
func (_m *mockLogSink) Init(info logr.RuntimeInfo) {
	_m.Called(info)
}

// mockLogSink_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type mockLogSink_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - info logr.RuntimeInfo
func (_e *mockLogSink_Expecter) Init(info interface{}) *mockLogSink_Init_Call {
	return &mockLogSink_Init_Call{Call: _e.mock.On("Init", info)}
}

func (_c *mockLogSink_Init_Call) Run(run func(info logr.RuntimeInfo)) *mockLogSink_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(logr.RuntimeInfo))
	})
	return _c
}

func (_c *mockLogSink_Init_Call) Return() *mockLogSink_Init_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockLogSink_Init_Call) RunAndReturn(run func(logr.RuntimeInfo)) *mockLogSink_Init_Call {
	_c.Run(run)
	return _c
}

// WithName provides a mock function with given fields: name
func (_m *mockLogSink) WithName(name string) logr.LogSink {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for WithName")
	}

	var r0 logr.LogSink
	if rf, ok := ret.Get(0).(func(string) logr.LogSink); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(logr.LogSink)
		}
	}

	return r0
}

// mockLogSink_WithName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithName'
type mockLogSink_WithName_Call struct {
	*mock.Call
}

// WithName is a helper method to define mock.On call
//   - name string
func (_e *mockLogSink_Expecter) WithName(name interface{}) *mockLogSink_WithName_Call {
	return &mockLogSink_WithName_Call{Call: _e.mock.On("WithName", name)}
}

func (_c *mockLogSink_WithName_Call) Run(run func(name string)) *mockLogSink_WithName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockLogSink_WithName_Call) Return(_a0 logr.LogSink) *mockLogSink_WithName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockLogSink_WithName_Call) RunAndReturn(run func(string) logr.LogSink) *mockLogSink_WithName_Call {
	_c.Call.Return(run)
	return _c
}

// WithValues provides a mock function with given fields: keysAndValues
func (_m *mockLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	var _ca []interface{}
	_ca = append(_ca, keysAndValues...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for WithValues")
	}

	var r0 logr.LogSink
	if rf, ok := ret.Get(0).(func(...interface{}) logr.LogSink); ok {
		r0 = rf(keysAndValues...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(logr.LogSink)
		}
	}

	return r0
}

// mockLogSink_WithValues_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithValues'
type mockLogSink_WithValues_Call struct {
	*mock.Call
}

// WithValues is a helper method to define mock.On call
//   - keysAndValues ...interface{}
func (_e *mockLogSink_Expecter) WithValues(keysAndValues ...interface{}) *mockLogSink_WithValues_Call {
	return &mockLogSink_WithValues_Call{Call: _e.mock.On("WithValues",
		append([]interface{}{}, keysAndValues...)...)}
}

func (_c *mockLogSink_WithValues_Call) Run(run func(keysAndValues ...interface{})) *mockLogSink_WithValues_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *mockLogSink_WithValues_Call) Return(_a0 logr.LogSink) *mockLogSink_WithValues_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockLogSink_WithValues_Call) RunAndReturn(run func(...interface{}) logr.LogSink) *mockLogSink_WithValues_Call {
	_c.Call.Return(run)
	return _c
}

// newMockLogSink creates a new instance of mockLogSink. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockLogSink(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockLogSink {
	mock := &mockLogSink{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
