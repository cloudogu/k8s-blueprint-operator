// Code generated by mockery v2.20.0. DO NOT EDIT.

package application

import (
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

// GetAll provides a mock function with given fields:
func (_m *mockDoguInstallationRepository) GetAll() ([]ecosystem.DoguInstallation, error) {
	ret := _m.Called()

	var r0 []ecosystem.DoguInstallation
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]ecosystem.DoguInstallation, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []ecosystem.DoguInstallation); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ecosystem.DoguInstallation)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
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
func (_e *mockDoguInstallationRepository_Expecter) GetAll() *mockDoguInstallationRepository_GetAll_Call {
	return &mockDoguInstallationRepository_GetAll_Call{Call: _e.mock.On("GetAll")}
}

func (_c *mockDoguInstallationRepository_GetAll_Call) Run(run func()) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockDoguInstallationRepository_GetAll_Call) Return(_a0 []ecosystem.DoguInstallation, _a1 error) *mockDoguInstallationRepository_GetAll_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockDoguInstallationRepository_GetAll_Call) RunAndReturn(run func() ([]ecosystem.DoguInstallation, error)) *mockDoguInstallationRepository_GetAll_Call {
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
