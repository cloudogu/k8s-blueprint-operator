// Code generated by mockery v2.20.0. DO NOT EDIT.

package maintenance

import (
	domainservice "github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	mock "github.com/stretchr/testify/mock"
)

// mockSwitcher is an autogenerated mock type for the switcher type
type mockSwitcher struct {
	mock.Mock
}

type mockSwitcher_Expecter struct {
	mock *mock.Mock
}

func (_m *mockSwitcher) EXPECT() *mockSwitcher_Expecter {
	return &mockSwitcher_Expecter{mock: &_m.Mock}
}

// activate provides a mock function with given fields: content
func (_m *mockSwitcher) activate(content domainservice.MaintenancePageModel) error {
	ret := _m.Called(content)

	var r0 error
	if rf, ok := ret.Get(0).(func(domainservice.MaintenancePageModel) error); ok {
		r0 = rf(content)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSwitcher_activate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'activate'
type mockSwitcher_activate_Call struct {
	*mock.Call
}

// activate is a helper method to define mock.On call
//   - content domainservice.MaintenancePageModel
func (_e *mockSwitcher_Expecter) activate(content interface{}) *mockSwitcher_activate_Call {
	return &mockSwitcher_activate_Call{Call: _e.mock.On("activate", content)}
}

func (_c *mockSwitcher_activate_Call) Run(run func(content domainservice.MaintenancePageModel)) *mockSwitcher_activate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domainservice.MaintenancePageModel))
	})
	return _c
}

func (_c *mockSwitcher_activate_Call) Return(_a0 error) *mockSwitcher_activate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSwitcher_activate_Call) RunAndReturn(run func(domainservice.MaintenancePageModel) error) *mockSwitcher_activate_Call {
	_c.Call.Return(run)
	return _c
}

// deactivate provides a mock function with given fields:
func (_m *mockSwitcher) deactivate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockSwitcher_deactivate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'deactivate'
type mockSwitcher_deactivate_Call struct {
	*mock.Call
}

// deactivate is a helper method to define mock.On call
func (_e *mockSwitcher_Expecter) deactivate() *mockSwitcher_deactivate_Call {
	return &mockSwitcher_deactivate_Call{Call: _e.mock.On("deactivate")}
}

func (_c *mockSwitcher_deactivate_Call) Run(run func()) *mockSwitcher_deactivate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockSwitcher_deactivate_Call) Return(_a0 error) *mockSwitcher_deactivate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockSwitcher_deactivate_Call) RunAndReturn(run func() error) *mockSwitcher_deactivate_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockSwitcher interface {
	mock.TestingT
	Cleanup(func())
}

// newMockSwitcher creates a new instance of mockSwitcher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockSwitcher(t mockConstructorTestingTnewMockSwitcher) *mockSwitcher {
	mock := &mockSwitcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}