// Code generated by mockery v2.53.3. DO NOT EDIT.

package kubernetes

import mock "github.com/stretchr/testify/mock"

// mockBlueprintGetter is an autogenerated mock type for the blueprintGetter type
type mockBlueprintGetter struct {
	mock.Mock
}

type mockBlueprintGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockBlueprintGetter) EXPECT() *mockBlueprintGetter_Expecter {
	return &mockBlueprintGetter_Expecter{mock: &_m.Mock}
}

// Blueprints provides a mock function with given fields: namespace
func (_m *mockBlueprintGetter) Blueprints(namespace string) BlueprintInterface {
	ret := _m.Called(namespace)

	if len(ret) == 0 {
		panic("no return value specified for Blueprints")
	}

	var r0 BlueprintInterface
	if rf, ok := ret.Get(0).(func(string) BlueprintInterface); ok {
		r0 = rf(namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(BlueprintInterface)
		}
	}

	return r0
}

// mockBlueprintGetter_Blueprints_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Blueprints'
type mockBlueprintGetter_Blueprints_Call struct {
	*mock.Call
}

// Blueprints is a helper method to define mock.On call
//   - namespace string
func (_e *mockBlueprintGetter_Expecter) Blueprints(namespace interface{}) *mockBlueprintGetter_Blueprints_Call {
	return &mockBlueprintGetter_Blueprints_Call{Call: _e.mock.On("Blueprints", namespace)}
}

func (_c *mockBlueprintGetter_Blueprints_Call) Run(run func(namespace string)) *mockBlueprintGetter_Blueprints_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *mockBlueprintGetter_Blueprints_Call) Return(_a0 BlueprintInterface) *mockBlueprintGetter_Blueprints_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockBlueprintGetter_Blueprints_Call) RunAndReturn(run func(string) BlueprintInterface) *mockBlueprintGetter_Blueprints_Call {
	_c.Call.Return(run)
	return _c
}

// newMockBlueprintGetter creates a new instance of mockBlueprintGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockBlueprintGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockBlueprintGetter {
	mock := &mockBlueprintGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
