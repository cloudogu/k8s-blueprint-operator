// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	dogu "github.com/cloudogu/ces-commons-lib/dogu"
	ecosystem "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"

	mock "github.com/stretchr/testify/mock"
)

// mockDoguInstallationRepository is an autogenerated mock type for the doguInstallationRepository type
type mockDoguInstallationRepository struct {
	mock.Mock
}

type mockDoguInstallationRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockDoguInstallationRepository) EXPECT() *mockDoguInstallationRepository_Expecter {
	return &mockDoguInstallationRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, _a1
func (_m *mockDoguInstallationRepository) Create(ctx context.Context, _a1 *ecosystem.DoguInstallation) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.DoguInstallation) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguInstallationRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockDoguInstallationRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *ecosystem.DoguInstallation
func (_e *mockDoguInstallationRepository_Expecter) Create(ctx interface{}, _a1 interface{}) *mockDoguInstallationRepository_Create_Call {
	return &mockDoguInstallationRepository_Create_Call{Call: _e.mock.On("Create", ctx, _a1)}
}

func (_c *mockDoguInstallationRepository_Create_Call) Run(run func(ctx context.Context, _a1 *ecosystem.DoguInstallation)) *mockDoguInstallationRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.DoguInstallation))
	})
	return _c
}

func (_c *mockDoguInstallationRepository_Create_Call) Return(_a0 error) *mockDoguInstallationRepository_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguInstallationRepository_Create_Call) RunAndReturn(run func(context.Context, *ecosystem.DoguInstallation) error) *mockDoguInstallationRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, doguName
func (_m *mockDoguInstallationRepository) Delete(ctx context.Context, doguName dogu.SimpleName) error {
	ret := _m.Called(ctx, doguName)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleName) error); ok {
		r0 = rf(ctx, doguName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguInstallationRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type mockDoguInstallationRepository_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName dogu.SimpleName
func (_e *mockDoguInstallationRepository_Expecter) Delete(ctx interface{}, doguName interface{}) *mockDoguInstallationRepository_Delete_Call {
	return &mockDoguInstallationRepository_Delete_Call{Call: _e.mock.On("Delete", ctx, doguName)}
}

func (_c *mockDoguInstallationRepository_Delete_Call) Run(run func(ctx context.Context, doguName dogu.SimpleName)) *mockDoguInstallationRepository_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.SimpleName))
	})
	return _c
}

func (_c *mockDoguInstallationRepository_Delete_Call) Return(_a0 error) *mockDoguInstallationRepository_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguInstallationRepository_Delete_Call) RunAndReturn(run func(context.Context, dogu.SimpleName) error) *mockDoguInstallationRepository_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function with given fields: ctx
func (_m *mockDoguInstallationRepository) GetAll(ctx context.Context) (map[dogu.SimpleName]*ecosystem.DoguInstallation, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 map[dogu.SimpleName]*ecosystem.DoguInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[dogu.SimpleName]*ecosystem.DoguInstallation, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[dogu.SimpleName]*ecosystem.DoguInstallation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[dogu.SimpleName]*ecosystem.DoguInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguInstallationRepository_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type mockDoguInstallationRepository_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockDoguInstallationRepository_Expecter) GetAll(ctx interface{}) *mockDoguInstallationRepository_GetAll_Call {
	return &mockDoguInstallationRepository_GetAll_Call{Call: _e.mock.On("GetAll", ctx)}
}

func (_c *mockDoguInstallationRepository_GetAll_Call) Run(run func(ctx context.Context)) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockDoguInstallationRepository_GetAll_Call) Return(_a0 map[dogu.SimpleName]*ecosystem.DoguInstallation, _a1 error) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguInstallationRepository_GetAll_Call) RunAndReturn(run func(context.Context) (map[dogu.SimpleName]*ecosystem.DoguInstallation, error)) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByName provides a mock function with given fields: ctx, doguName
func (_m *mockDoguInstallationRepository) GetByName(ctx context.Context, doguName dogu.SimpleName) (*ecosystem.DoguInstallation, error) {
	ret := _m.Called(ctx, doguName)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

	var r0 *ecosystem.DoguInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleName) (*ecosystem.DoguInstallation, error)); ok {
		return rf(ctx, doguName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dogu.SimpleName) *ecosystem.DoguInstallation); ok {
		r0 = rf(ctx, doguName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.DoguInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dogu.SimpleName) error); ok {
		r1 = rf(ctx, doguName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockDoguInstallationRepository_GetByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByName'
type mockDoguInstallationRepository_GetByName_Call struct {
	*mock.Call
}

// GetByName is a helper method to define mock.On call
//   - ctx context.Context
//   - doguName dogu.SimpleName
func (_e *mockDoguInstallationRepository_Expecter) GetByName(ctx interface{}, doguName interface{}) *mockDoguInstallationRepository_GetByName_Call {
	return &mockDoguInstallationRepository_GetByName_Call{Call: _e.mock.On("GetByName", ctx, doguName)}
}

func (_c *mockDoguInstallationRepository_GetByName_Call) Run(run func(ctx context.Context, doguName dogu.SimpleName)) *mockDoguInstallationRepository_GetByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.SimpleName))
	})
	return _c
}

func (_c *mockDoguInstallationRepository_GetByName_Call) Return(_a0 *ecosystem.DoguInstallation, _a1 error) *mockDoguInstallationRepository_GetByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguInstallationRepository_GetByName_Call) RunAndReturn(run func(context.Context, dogu.SimpleName) (*ecosystem.DoguInstallation, error)) *mockDoguInstallationRepository_GetByName_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *mockDoguInstallationRepository) Update(ctx context.Context, _a1 *ecosystem.DoguInstallation) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *ecosystem.DoguInstallation) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockDoguInstallationRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockDoguInstallationRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 *ecosystem.DoguInstallation
func (_e *mockDoguInstallationRepository_Expecter) Update(ctx interface{}, _a1 interface{}) *mockDoguInstallationRepository_Update_Call {
	return &mockDoguInstallationRepository_Update_Call{Call: _e.mock.On("Update", ctx, _a1)}
}

func (_c *mockDoguInstallationRepository_Update_Call) Run(run func(ctx context.Context, _a1 *ecosystem.DoguInstallation)) *mockDoguInstallationRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*ecosystem.DoguInstallation))
	})
	return _c
}

func (_c *mockDoguInstallationRepository_Update_Call) Return(_a0 error) *mockDoguInstallationRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockDoguInstallationRepository_Update_Call) RunAndReturn(run func(context.Context, *ecosystem.DoguInstallation) error) *mockDoguInstallationRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// newMockDoguInstallationRepository creates a new instance of mockDoguInstallationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockDoguInstallationRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockDoguInstallationRepository {
	mock := &mockDoguInstallationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
