// Code generated by mockery v2.53.3. DO NOT EDIT.

package kubernetes

import (
	mock "github.com/stretchr/testify/mock"
	flowcontrol "k8s.io/client-go/util/flowcontrol"

	rest "k8s.io/client-go/rest"

	schema "k8s.io/apimachinery/pkg/runtime/schema"

	types "k8s.io/apimachinery/pkg/types"
)

// mockRestInterface is an autogenerated mock type for the restInterface type
type mockRestInterface struct {
	mock.Mock
}

type mockRestInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *mockRestInterface) EXPECT() *mockRestInterface_Expecter {
	return &mockRestInterface_Expecter{mock: &_m.Mock}
}

// APIVersion provides a mock function with no fields
func (_m *mockRestInterface) APIVersion() schema.GroupVersion {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for APIVersion")
	}

	var r0 schema.GroupVersion
	if rf, ok := ret.Get(0).(func() schema.GroupVersion); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(schema.GroupVersion)
	}

	return r0
}

// mockRestInterface_APIVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'APIVersion'
type mockRestInterface_APIVersion_Call struct {
	*mock.Call
}

// APIVersion is a helper method to define mock.On call
func (_e *mockRestInterface_Expecter) APIVersion() *mockRestInterface_APIVersion_Call {
	return &mockRestInterface_APIVersion_Call{Call: _e.mock.On("APIVersion")}
}

func (_c *mockRestInterface_APIVersion_Call) Run(run func()) *mockRestInterface_APIVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRestInterface_APIVersion_Call) Return(_a0 schema.GroupVersion) *mockRestInterface_APIVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_APIVersion_Call) RunAndReturn(run func() schema.GroupVersion) *mockRestInterface_APIVersion_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with no fields
func (_m *mockRestInterface) Delete() *rest.Request {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func() *rest.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// mockRestInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockRestInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
func (_e *mockRestInterface_Expecter) Delete() *mockRestInterface_Delete_Call {
	return &mockRestInterface_Delete_Call{Call: _e.mock.On("Delete")}
}

func (_c *mockRestInterface_Delete_Call) Run(run func()) *mockRestInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRestInterface_Delete_Call) Return(_a0 *rest.Request) *mockRestInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_Delete_Call) RunAndReturn(run func() *rest.Request) *mockRestInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with no fields
func (_m *mockRestInterface) Get() *rest.Request {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func() *rest.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// mockRestInterface_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockRestInterface_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
func (_e *mockRestInterface_Expecter) Get() *mockRestInterface_Get_Call {
	return &mockRestInterface_Get_Call{Call: _e.mock.On("Get")}
}

func (_c *mockRestInterface_Get_Call) Run(run func()) *mockRestInterface_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRestInterface_Get_Call) Return(_a0 *rest.Request) *mockRestInterface_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_Get_Call) RunAndReturn(run func() *rest.Request) *mockRestInterface_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetRateLimiter provides a mock function with no fields
func (_m *mockRestInterface) GetRateLimiter() flowcontrol.RateLimiter {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRateLimiter")
	}

	var r0 flowcontrol.RateLimiter
	if rf, ok := ret.Get(0).(func() flowcontrol.RateLimiter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(flowcontrol.RateLimiter)
		}
	}

	return r0
}

// mockRestInterface_GetRateLimiter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRateLimiter'
type mockRestInterface_GetRateLimiter_Call struct {
	*mock.Call
}

// GetRateLimiter is a helper method to define mock.On call
func (_e *mockRestInterface_Expecter) GetRateLimiter() *mockRestInterface_GetRateLimiter_Call {
	return &mockRestInterface_GetRateLimiter_Call{Call: _e.mock.On("GetRateLimiter")}
}

func (_c *mockRestInterface_GetRateLimiter_Call) Run(run func()) *mockRestInterface_GetRateLimiter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRestInterface_GetRateLimiter_Call) Return(_a0 flowcontrol.RateLimiter) *mockRestInterface_GetRateLimiter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_GetRateLimiter_Call) RunAndReturn(run func() flowcontrol.RateLimiter) *mockRestInterface_GetRateLimiter_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: pt
func (_m *mockRestInterface) Patch(pt types.PatchType) *rest.Request {
	ret := _m.Called(pt)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func(types.PatchType) *rest.Request); ok {
		r0 = rf(pt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// mockRestInterface_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type mockRestInterface_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - pt types.PatchType
func (_e *mockRestInterface_Expecter) Patch(pt interface{}) *mockRestInterface_Patch_Call {
	return &mockRestInterface_Patch_Call{Call: _e.mock.On("Patch", pt)}
}

func (_c *mockRestInterface_Patch_Call) Run(run func(pt types.PatchType)) *mockRestInterface_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.PatchType))
	})
	return _c
}

func (_c *mockRestInterface_Patch_Call) Return(_a0 *rest.Request) *mockRestInterface_Patch_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_Patch_Call) RunAndReturn(run func(types.PatchType) *rest.Request) *mockRestInterface_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Post provides a mock function with no fields
func (_m *mockRestInterface) Post() *rest.Request {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Post")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func() *rest.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// mockRestInterface_Post_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Post'
type mockRestInterface_Post_Call struct {
	*mock.Call
}

// Post is a helper method to define mock.On call
func (_e *mockRestInterface_Expecter) Post() *mockRestInterface_Post_Call {
	return &mockRestInterface_Post_Call{Call: _e.mock.On("Post")}
}

func (_c *mockRestInterface_Post_Call) Run(run func()) *mockRestInterface_Post_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRestInterface_Post_Call) Return(_a0 *rest.Request) *mockRestInterface_Post_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_Post_Call) RunAndReturn(run func() *rest.Request) *mockRestInterface_Post_Call {
	_c.Call.Return(run)
	return _c
}

// Put provides a mock function with no fields
func (_m *mockRestInterface) Put() *rest.Request {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Put")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func() *rest.Request); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// mockRestInterface_Put_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Put'
type mockRestInterface_Put_Call struct {
	*mock.Call
}

// Put is a helper method to define mock.On call
func (_e *mockRestInterface_Expecter) Put() *mockRestInterface_Put_Call {
	return &mockRestInterface_Put_Call{Call: _e.mock.On("Put")}
}

func (_c *mockRestInterface_Put_Call) Run(run func()) *mockRestInterface_Put_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockRestInterface_Put_Call) Return(_a0 *rest.Request) *mockRestInterface_Put_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_Put_Call) RunAndReturn(run func() *rest.Request) *mockRestInterface_Put_Call {
	_c.Call.Return(run)
	return _c
}

// Verb provides a mock function with given fields: verb
func (_m *mockRestInterface) Verb(verb string) *rest.Request {
	ret := _m.Called(verb)

	if len(ret) == 0 {
		panic("no return value specified for Verb")
	}

	var r0 *rest.Request
	if rf, ok := ret.Get(0).(func(string) *rest.Request); ok {
		r0 = rf(verb)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rest.Request)
		}
	}

	return r0
}

// mockRestInterface_Verb_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Verb'
type mockRestInterface_Verb_Call struct {
	*mock.Call
}

// Verb is a helper method to define mock.On call
//   - verb string
func (_e *mockRestInterface_Expecter) Verb(verb interface{}) *mockRestInterface_Verb_Call {
	return &mockRestInterface_Verb_Call{Call: _e.mock.On("Verb", verb)}
}

func (_c *mockRestInterface_Verb_Call) Run(run func(verb string)) *mockRestInterface_Verb_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockRestInterface_Verb_Call) Return(_a0 *rest.Request) *mockRestInterface_Verb_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockRestInterface_Verb_Call) RunAndReturn(run func(string) *rest.Request) *mockRestInterface_Verb_Call {
	_c.Call.Return(run)
	return _c
}

// newMockRestInterface creates a new instance of mockRestInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockRestInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockRestInterface {
	mock := &mockRestInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
