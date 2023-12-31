// Code generated by mockery v2.20.0. DO NOT EDIT.

package application

import (
	context "context"

	ecosystem "github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
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

// GetAll provides a mock function with given fields: ctx
func (_m *mockDoguInstallationRepository) GetAll(ctx context.Context) (map[string]*ecosystem.DoguInstallation, error) {
	ret := _m.Called(ctx)

	var r0 map[string]*ecosystem.DoguInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (map[string]*ecosystem.DoguInstallation, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) map[string]*ecosystem.DoguInstallation); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*ecosystem.DoguInstallation)
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

func (_c *mockDoguInstallationRepository_GetAll_Call) Return(_a0 map[string]*ecosystem.DoguInstallation, _a1 error) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguInstallationRepository_GetAll_Call) RunAndReturn(run func(context.Context) (map[string]*ecosystem.DoguInstallation, error)) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// GetByName provides a mock function with given fields: ctx, doguName
func (_m *mockDoguInstallationRepository) GetByName(ctx context.Context, doguName string) (*ecosystem.DoguInstallation, error) {
	ret := _m.Called(ctx, doguName)

	var r0 *ecosystem.DoguInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*ecosystem.DoguInstallation, error)); ok {
		return rf(ctx, doguName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *ecosystem.DoguInstallation); ok {
		r0 = rf(ctx, doguName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ecosystem.DoguInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
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
//   - doguName string
func (_e *mockDoguInstallationRepository_Expecter) GetByName(ctx interface{}, doguName interface{}) *mockDoguInstallationRepository_GetByName_Call {
	return &mockDoguInstallationRepository_GetByName_Call{Call: _e.mock.On("GetByName", ctx, doguName)}
}

func (_c *mockDoguInstallationRepository_GetByName_Call) Run(run func(ctx context.Context, doguName string)) *mockDoguInstallationRepository_GetByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockDoguInstallationRepository_GetByName_Call) Return(_a0 *ecosystem.DoguInstallation, _a1 error) *mockDoguInstallationRepository_GetByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguInstallationRepository_GetByName_Call) RunAndReturn(run func(context.Context, string) (*ecosystem.DoguInstallation, error)) *mockDoguInstallationRepository_GetByName_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockDoguInstallationRepository interface {
	mock.TestingT
	Cleanup(func())
}

// newMockDoguInstallationRepository creates a new instance of mockDoguInstallationRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockDoguInstallationRepository(t mockConstructorTestingTnewMockDoguInstallationRepository) *mockDoguInstallationRepository {
	mock := &mockDoguInstallationRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
