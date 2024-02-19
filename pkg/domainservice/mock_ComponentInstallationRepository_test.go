// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// MockComponentInstallationRepository is an autogenerated mock type for the ComponentInstallationRepository type
type MockComponentInstallationRepository struct {
	mock.Mock
}

type MockComponentInstallationRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockComponentInstallationRepository) EXPECT() *MockComponentInstallationRepository_Expecter {
	return &MockComponentInstallationRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, component
func (_m *MockComponentInstallationRepository) Create(ctx context.Context, component *ecosystem.ComponentInstallation) error {
	ret := _m.Called(ctx, component)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.ComponentInstallation) error); ok {
		r0 = rf(ctx, component)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockComponentInstallationRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockComponentInstallationRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - component *ecosystem.ComponentInstallation
func (_e *MockComponentInstallationRepository_Expecter) Create(ctx interface{}, component interface{}) *MockComponentInstallationRepository_Create_Call {
	return &MockComponentInstallationRepository_Create_Call{Call: _e.mock.On("Create", ctx, component)}
}

func (_c *MockComponentInstallationRepository_Create_Call) Run(run func(ctx context.Context, component *ecosystem.ComponentInstallation)) *MockComponentInstallationRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.ComponentInstallation))
	})
	return _c
}

func (_c *MockComponentInstallationRepository_Create_Call) Return(_a0 error) *MockComponentInstallationRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockComponentInstallationRepository_Create_Call) RunAndReturn(run func(context.Context, *ecosystem.ComponentInstallation) error) *MockComponentInstallationRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, componentName
func (_m *MockComponentInstallationRepository) Delete(ctx context.Context, componentName common.SimpleComponentName) error {
	ret := _m.Called(ctx, componentName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleComponentName) error); ok {
		r0 = rf(ctx, componentName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockComponentInstallationRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockComponentInstallationRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName common.SimpleComponentName
func (_e *MockComponentInstallationRepository_Expecter) Delete(ctx interface{}, componentName interface{}) *MockComponentInstallationRepository_Delete_Call {
	return &MockComponentInstallationRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, componentName)}
}

func (_c *MockComponentInstallationRepository_Delete_Call) Run(run func(ctx context.Context, componentName common.SimpleComponentName)) *MockComponentInstallationRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleComponentName))
	})
	return _c
}

func (_c *MockComponentInstallationRepository_Delete_Call) Return(_a0 error) *MockComponentInstallationRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockComponentInstallationRepository_Delete_Call) RunAndReturn(run func(context.Context, common.SimpleComponentName) error) *MockComponentInstallationRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *MockComponentInstallationRepository) GetAll(ctx context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error) {
	ret := _m.Called(ctx)

	var r0 map[common.SimpleComponentName]*ecosystem.ComponentInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[common.SimpleComponentName]*ecosystem.ComponentInstallation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[common.SimpleComponentName]*ecosystem.ComponentInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockComponentInstallationRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockComponentInstallationRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockComponentInstallationRepository_Expecter) GetAll(ctx interface{}) *MockComponentInstallationRepository_GetAll_Call {
	return &MockComponentInstallationRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *MockComponentInstallationRepository_GetAll_Call) Run(run func(ctx context.Context)) *MockComponentInstallationRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockComponentInstallationRepository_GetAll_Call) Return(_a0 map[common.SimpleComponentName]*ecosystem.ComponentInstallation, _a1 error) *MockComponentInstallationRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockComponentInstallationRepository_GetAll_Call) RunAndReturn(run func(context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error)) *MockComponentInstallationRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByName provides a mock function with given fields: ctx, componentName
func (_m *MockComponentInstallationRepository) GetByName(ctx context.Context, componentName common.SimpleComponentName) (*ecosystem.ComponentInstallation, error) {
	ret := _m.Called(ctx, componentName)

	var r0 *ecosystem.ComponentInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleComponentName) (*ecosystem.ComponentInstallation, error)); ok {
		return rf(ctx, componentName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleComponentName) *ecosystem.ComponentInstallation); ok {
		r0 = rf(ctx, componentName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.ComponentInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.SimpleComponentName) error); ok {
		r1 = rf(ctx, componentName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockComponentInstallationRepository_GetByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByName'
type MockComponentInstallationRepository_GetByName_Call struct {
	*mock.Call
}

// GetByName is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName common.SimpleComponentName
func (_e *MockComponentInstallationRepository_Expecter) GetByName(ctx interface{}, componentName interface{}) *MockComponentInstallationRepository_GetByName_Call {
	return &MockComponentInstallationRepository_GetByName_Call{Call: _e.mock.On("GetByName", ctx, componentName)}
}

func (_c *MockComponentInstallationRepository_GetByName_Call) Run(run func(ctx context.Context, componentName common.SimpleComponentName)) *MockComponentInstallationRepository_GetByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleComponentName))
	})
	return _c
}

func (_c *MockComponentInstallationRepository_GetByName_Call) Return(_a0 *ecosystem.ComponentInstallation, _a1 error) *MockComponentInstallationRepository_GetByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockComponentInstallationRepository_GetByName_Call) RunAndReturn(run func(context.Context, common.SimpleComponentName) (*ecosystem.ComponentInstallation, error)) *MockComponentInstallationRepository_GetByName_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, component
func (_m *MockComponentInstallationRepository) Update(ctx context.Context, component *ecosystem.ComponentInstallation) error {
	ret := _m.Called(ctx, component)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.ComponentInstallation) error); ok {
		r0 = rf(ctx, component)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockComponentInstallationRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockComponentInstallationRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - component *ecosystem.ComponentInstallation
func (_e *MockComponentInstallationRepository_Expecter) Update(ctx interface{}, component interface{}) *MockComponentInstallationRepository_Update_Call {
	return &MockComponentInstallationRepository_Update_Call{Call: _e.mock.On("Update", ctx, component)}
}

func (_c *MockComponentInstallationRepository_Update_Call) Run(run func(ctx context.Context, component *ecosystem.ComponentInstallation)) *MockComponentInstallationRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.ComponentInstallation))
	})
	return _c
}

func (_c *MockComponentInstallationRepository_Update_Call) Return(_a0 error) *MockComponentInstallationRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockComponentInstallationRepository_Update_Call) RunAndReturn(run func(context.Context, *ecosystem.ComponentInstallation) error) *MockComponentInstallationRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockComponentInstallationRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockComponentInstallationRepository creates a new instance of MockComponentInstallationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockComponentInstallationRepository(t mockConstructorTestingTNewMockComponentInstallationRepository) *MockComponentInstallationRepository {
	mock := &MockComponentInstallationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
