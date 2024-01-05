// Code generated by mockery v2.20.0. DO NOT EDIT.

package serializer

import (
	domain "github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	mock "github.com/stretchr/testify/mock"
)

// MockBlueprintSerializer is an autogenerated mock type for the BlueprintSerializer type
type MockBlueprintSerializer struct {
	mock.Mock
}

type MockBlueprintSerializer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBlueprintSerializer) EXPECT() *MockBlueprintSerializer_Expecter {
	return &MockBlueprintSerializer_Expecter{mock: &_m.Mock}
}

// Deserialize provides a mock function with given fields: rawBlueprint
func (_m *MockBlueprintSerializer) Deserialize(rawBlueprint string) (domain.Blueprint, error) {
	ret := _m.Called(rawBlueprint)

	var r0 domain.Blueprint
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.Blueprint, error)); ok {
		return rf(rawBlueprint)
	}
	if rf, ok := ret.Get(0).(func(string) domain.Blueprint); ok {
		r0 = rf(rawBlueprint)
	} else {
		r0 = ret.Get(0).(domain.Blueprint)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(rawBlueprint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBlueprintSerializer_Deserialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deserialize'
type MockBlueprintSerializer_Deserialize_Call struct {
	*mock.Call
}

// Deserialize is a helper method to define mock.On call
//   - rawBlueprint string
func (_e *MockBlueprintSerializer_Expecter) Deserialize(rawBlueprint interface{}) *MockBlueprintSerializer_Deserialize_Call {
	return &MockBlueprintSerializer_Deserialize_Call{Call: _e.mock.On("Deserialize", rawBlueprint)}
}

func (_c *MockBlueprintSerializer_Deserialize_Call) Run(run func(rawBlueprint string)) *MockBlueprintSerializer_Deserialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockBlueprintSerializer_Deserialize_Call) Return(_a0 domain.Blueprint, _a1 error) *MockBlueprintSerializer_Deserialize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBlueprintSerializer_Deserialize_Call) RunAndReturn(run func(string) (domain.Blueprint, error)) *MockBlueprintSerializer_Deserialize_Call {
	_c.Call.Return(run)
	return _c
}

// Serialize provides a mock function with given fields: blueprint
func (_m *MockBlueprintSerializer) Serialize(blueprint domain.Blueprint) (string, error) {
	ret := _m.Called(blueprint)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Blueprint) (string, error)); ok {
		return rf(blueprint)
	}
	if rf, ok := ret.Get(0).(func(domain.Blueprint) string); ok {
		r0 = rf(blueprint)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(domain.Blueprint) error); ok {
		r1 = rf(blueprint)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockBlueprintSerializer_Serialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serialize'
type MockBlueprintSerializer_Serialize_Call struct {
	*mock.Call
}

// Serialize is a helper method to define mock.On call
//   - blueprint domain.Blueprint
func (_e *MockBlueprintSerializer_Expecter) Serialize(blueprint interface{}) *MockBlueprintSerializer_Serialize_Call {
	return &MockBlueprintSerializer_Serialize_Call{Call: _e.mock.On("Serialize", blueprint)}
}

func (_c *MockBlueprintSerializer_Serialize_Call) Run(run func(blueprint domain.Blueprint)) *MockBlueprintSerializer_Serialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(domain.Blueprint))
	})
	return _c
}

func (_c *MockBlueprintSerializer_Serialize_Call) Return(_a0 string, _a1 error) *MockBlueprintSerializer_Serialize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockBlueprintSerializer_Serialize_Call) RunAndReturn(run func(domain.Blueprint) (string, error)) *MockBlueprintSerializer_Serialize_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockBlueprintSerializer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockBlueprintSerializer creates a new instance of MockBlueprintSerializer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockBlueprintSerializer(t mockConstructorTestingTNewMockBlueprintSerializer) *MockBlueprintSerializer {
	mock := &MockBlueprintSerializer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
