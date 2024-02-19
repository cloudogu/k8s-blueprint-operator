// Code generated by mockery v2.20.0. DO NOT EDIT.

package domain

import mock "github.com/stretchr/testify/mock"

// MockEvent is an autogenerated mock type for the Event type
type MockEvent struct {
	mock.Mock
}

type MockEvent_Expecter struct {
	mock *mock.Mock
}

func (_m *MockEvent) EXPECT() *MockEvent_Expecter {
	return &MockEvent_Expecter{mock: &_m.Mock}
}

// Message provides a mock function with given fields:
func (_m *MockEvent) Message() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockEvent_Message_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Message'
type MockEvent_Message_Call struct {
	*mock.Call
}

// Message is a helper method to define mock.On call
func (_e *MockEvent_Expecter) Message() *MockEvent_Message_Call {
	return &MockEvent_Message_Call{Call: _e.mock.On("Message")}
}

func (_c *MockEvent_Message_Call) Run(run func()) *MockEvent_Message_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEvent_Message_Call) Return(_a0 string) *MockEvent_Message_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEvent_Message_Call) RunAndReturn(run func() string) *MockEvent_Message_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with given fields:
func (_m *MockEvent) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockEvent_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SimpleName'
type MockEvent_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockEvent_Expecter) Name() *MockEvent_Name_Call {
	return &MockEvent_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockEvent_Name_Call) Run(run func()) *MockEvent_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockEvent_Name_Call) Return(_a0 string) *MockEvent_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockEvent_Name_Call) RunAndReturn(run func() string) *MockEvent_Name_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockEvent interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockEvent creates a new instance of MockEvent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockEvent(t mockConstructorTestingTNewMockEvent) *MockEvent {
	mock := &MockEvent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
