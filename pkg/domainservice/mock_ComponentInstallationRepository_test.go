// Code generated by mockery v2.20.0. DO NOT EDIT.

package domainservice

import (
	context "context"

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

// GetAll provides a mock function with given fields: ctx
func (_m *MockComponentInstallationRepository) GetAll(ctx context.Context) (map[string]*ecosystem.ComponentInstallation, error) {
	ret := _m.Called(ctx)

	var r0 map[string]*ecosystem.ComponentInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[string]*ecosystem.ComponentInstallation, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]*ecosystem.ComponentInstallation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*ecosystem.ComponentInstallation)
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

func (_c *MockComponentInstallationRepository_GetAll_Call) Return(_a0 map[string]*ecosystem.ComponentInstallation, _a1 error) *MockComponentInstallationRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockComponentInstallationRepository_GetAll_Call) RunAndReturn(run func(context.Context) (map[string]*ecosystem.ComponentInstallation, error)) *MockComponentInstallationRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByName provides a mock function with given fields: ctx, componentName
func (_m *MockComponentInstallationRepository) GetByName(ctx context.Context, componentName string) (*ecosystem.ComponentInstallation, error) {
	ret := _m.Called(ctx, componentName)

	var r0 *ecosystem.ComponentInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*ecosystem.ComponentInstallation, error)); ok {
		return rf(ctx, componentName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *ecosystem.ComponentInstallation); ok {
		r0 = rf(ctx, componentName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.ComponentInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
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
//   - componentName string
func (_e *MockComponentInstallationRepository_Expecter) GetByName(ctx interface{}, componentName interface{}) *MockComponentInstallationRepository_GetByName_Call {
	return &MockComponentInstallationRepository_GetByName_Call{Call: _e.mock.On("GetByName", ctx, componentName)}
}

func (_c *MockComponentInstallationRepository_GetByName_Call) Run(run func(ctx context.Context, componentName string)) *MockComponentInstallationRepository_GetByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockComponentInstallationRepository_GetByName_Call) Return(_a0 *ecosystem.ComponentInstallation, _a1 error) *MockComponentInstallationRepository_GetByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockComponentInstallationRepository_GetByName_Call) RunAndReturn(run func(context.Context, string) (*ecosystem.ComponentInstallation, error)) *MockComponentInstallationRepository_GetByName_Call {
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