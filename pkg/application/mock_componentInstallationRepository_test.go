// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	common "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// mockComponentInstallationRepository is an autogenerated mock type for the componentInstallationRepository type
type mockComponentInstallationRepository struct {
	mock.Mock
}

type mockComponentInstallationRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockComponentInstallationRepository) EXPECT() *mockComponentInstallationRepository_Expecter {
	return &mockComponentInstallationRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, component
func (_m *mockComponentInstallationRepository) Create(ctx context.Context, component *ecosystem.ComponentInstallation) error {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.ComponentInstallation) error); ok {
		r0 = rf(ctx, component)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInstallationRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockComponentInstallationRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - component *ecosystem.ComponentInstallation
func (_e *mockComponentInstallationRepository_Expecter) Create(ctx interface{}, component interface{}) *mockComponentInstallationRepository_Create_Call {
	return &mockComponentInstallationRepository_Create_Call{Call: _e.mock.On("Create", ctx, component)}
}

func (_c *mockComponentInstallationRepository_Create_Call) Run(run func(ctx context.Context, component *ecosystem.ComponentInstallation)) *mockComponentInstallationRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.ComponentInstallation))
	})
	return _c
}

func (_c *mockComponentInstallationRepository_Create_Call) Return(_a0 error) *mockComponentInstallationRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInstallationRepository_Create_Call) RunAndReturn(run func(context.Context, *ecosystem.ComponentInstallation) error) *mockComponentInstallationRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, componentName
func (_m *mockComponentInstallationRepository) Delete(ctx context.Context, componentName common.SimpleComponentName) error {
	ret := _m.Called(ctx, componentName)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.SimpleComponentName) error); ok {
		r0 = rf(ctx, componentName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInstallationRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockComponentInstallationRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName common.SimpleComponentName
func (_e *mockComponentInstallationRepository_Expecter) Delete(ctx interface{}, componentName interface{}) *mockComponentInstallationRepository_Delete_Call {
	return &mockComponentInstallationRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, componentName)}
}

func (_c *mockComponentInstallationRepository_Delete_Call) Run(run func(ctx context.Context, componentName common.SimpleComponentName)) *mockComponentInstallationRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleComponentName))
	})
	return _c
}

func (_c *mockComponentInstallationRepository_Delete_Call) Return(_a0 error) *mockComponentInstallationRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInstallationRepository_Delete_Call) RunAndReturn(run func(context.Context, common.SimpleComponentName) error) *mockComponentInstallationRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *mockComponentInstallationRepository) GetAll(ctx context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

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

// mockComponentInstallationRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type mockComponentInstallationRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockComponentInstallationRepository_Expecter) GetAll(ctx interface{}) *mockComponentInstallationRepository_GetAll_Call {
	return &mockComponentInstallationRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *mockComponentInstallationRepository_GetAll_Call) Run(run func(ctx context.Context)) *mockComponentInstallationRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockComponentInstallationRepository_GetAll_Call) Return(_a0 map[common.SimpleComponentName]*ecosystem.ComponentInstallation, _a1 error) *mockComponentInstallationRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInstallationRepository_GetAll_Call) RunAndReturn(run func(context.Context) (map[common.SimpleComponentName]*ecosystem.ComponentInstallation, error)) *mockComponentInstallationRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByName provides a mock function with given fields: ctx, componentName
func (_m *mockComponentInstallationRepository) GetByName(ctx context.Context, componentName common.SimpleComponentName) (*ecosystem.ComponentInstallation, error) {
	ret := _m.Called(ctx, componentName)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

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

// mockComponentInstallationRepository_GetByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByName'
type mockComponentInstallationRepository_GetByName_Call struct {
	*mock.Call
}

// GetByName is a helper method to define mock.On call
//   - ctx context.Context
//   - componentName common.SimpleComponentName
func (_e *mockComponentInstallationRepository_Expecter) GetByName(ctx interface{}, componentName interface{}) *mockComponentInstallationRepository_GetByName_Call {
	return &mockComponentInstallationRepository_GetByName_Call{Call: _e.mock.On("GetByName", ctx, componentName)}
}

func (_c *mockComponentInstallationRepository_GetByName_Call) Run(run func(ctx context.Context, componentName common.SimpleComponentName)) *mockComponentInstallationRepository_GetByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.SimpleComponentName))
	})
	return _c
}

func (_c *mockComponentInstallationRepository_GetByName_Call) Return(_a0 *ecosystem.ComponentInstallation, _a1 error) *mockComponentInstallationRepository_GetByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockComponentInstallationRepository_GetByName_Call) RunAndReturn(run func(context.Context, common.SimpleComponentName) (*ecosystem.ComponentInstallation, error)) *mockComponentInstallationRepository_GetByName_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, component
func (_m *mockComponentInstallationRepository) Update(ctx context.Context, component *ecosystem.ComponentInstallation) error {
	ret := _m.Called(ctx, component)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.ComponentInstallation) error); ok {
		r0 = rf(ctx, component)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockComponentInstallationRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockComponentInstallationRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - component *ecosystem.ComponentInstallation
func (_e *mockComponentInstallationRepository_Expecter) Update(ctx interface{}, component interface{}) *mockComponentInstallationRepository_Update_Call {
	return &mockComponentInstallationRepository_Update_Call{Call: _e.mock.On("Update", ctx, component)}
}

func (_c *mockComponentInstallationRepository_Update_Call) Run(run func(ctx context.Context, component *ecosystem.ComponentInstallation)) *mockComponentInstallationRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.ComponentInstallation))
	})
	return _c
}

func (_c *mockComponentInstallationRepository_Update_Call) Return(_a0 error) *mockComponentInstallationRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockComponentInstallationRepository_Update_Call) RunAndReturn(run func(context.Context, *ecosystem.ComponentInstallation) error) *mockComponentInstallationRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// newMockComponentInstallationRepository creates a new instance of mockComponentInstallationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockComponentInstallationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockComponentInstallationRepository {
	mock := &mockComponentInstallationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
