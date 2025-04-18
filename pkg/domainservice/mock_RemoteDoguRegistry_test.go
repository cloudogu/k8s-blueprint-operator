// Code generated by mockery v2.53.3. DO NOT EDIT.

package domainservice

import (
	context "context"

	dogu "github.com/cloudogu/ces-commons-lib/dogu"
	core "github.com/cloudogu/cesapp-lib/core"

	mock "github.com/stretchr/testify/mock"
)

// MockRemoteDoguRegistry is an autogenerated mock type for the RemoteDoguRegistry type
type MockRemoteDoguRegistry struct {
	mock.Mock
}

type MockRemoteDoguRegistry_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRemoteDoguRegistry) EXPECT() *MockRemoteDoguRegistry_Expecter {
	return &MockRemoteDoguRegistry_Expecter{mock: &_m.Mock}
}

// GetDogu provides a mock function with given fields: ctx, qualifiedDoguVersion
func (_m *MockRemoteDoguRegistry) GetDogu(ctx context.Context, qualifiedDoguVersion dogu.QualifiedVersion) (*core.Dogu, error) {
	ret := _m.Called(ctx, qualifiedDoguVersion)

	if len(ret) == 0 {
		panic("no return value specified for GetDogu")
	}

	var r0 *core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dogu.QualifiedVersion) (*core.Dogu, error)); ok {
		return rf(ctx, qualifiedDoguVersion)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dogu.QualifiedVersion) *core.Dogu); ok {
		r0 = rf(ctx, qualifiedDoguVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, dogu.QualifiedVersion) error); ok {
		r1 = rf(ctx, qualifiedDoguVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRemoteDoguRegistry_GetDogu_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDogu'
type MockRemoteDoguRegistry_GetDogu_Call struct {
	*mock.Call
}

// GetDogu is a helper method to define mock.On call
//   - ctx context.Context
//   - qualifiedDoguVersion dogu.QualifiedVersion
func (_e *MockRemoteDoguRegistry_Expecter) GetDogu(ctx interface{}, qualifiedDoguVersion interface{}) *MockRemoteDoguRegistry_GetDogu_Call {
	return &MockRemoteDoguRegistry_GetDogu_Call{Call: _e.mock.On("GetDogu", ctx, qualifiedDoguVersion)}
}

func (_c *MockRemoteDoguRegistry_GetDogu_Call) Run(run func(ctx context.Context, qualifiedDoguVersion dogu.QualifiedVersion)) *MockRemoteDoguRegistry_GetDogu_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dogu.QualifiedVersion))
	})
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogu_Call) Return(_a0 *core.Dogu, _a1 error) *MockRemoteDoguRegistry_GetDogu_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogu_Call) RunAndReturn(run func(context.Context, dogu.QualifiedVersion) (*core.Dogu, error)) *MockRemoteDoguRegistry_GetDogu_Call {
	_c.Call.Return(run)
	return _c
}

// GetDogus provides a mock function with given fields: ctx, dogusToLoad
func (_m *MockRemoteDoguRegistry) GetDogus(ctx context.Context, dogusToLoad []dogu.QualifiedVersion) (map[dogu.QualifiedName]*core.Dogu, error) {
	ret := _m.Called(ctx, dogusToLoad)

	if len(ret) == 0 {
		panic("no return value specified for GetDogus")
	}

	var r0 map[dogu.QualifiedName]*core.Dogu
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.QualifiedVersion) (map[dogu.QualifiedName]*core.Dogu, error)); ok {
		return rf(ctx, dogusToLoad)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []dogu.QualifiedVersion) map[dogu.QualifiedName]*core.Dogu); ok {
		r0 = rf(ctx, dogusToLoad)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[dogu.QualifiedName]*core.Dogu)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []dogu.QualifiedVersion) error); ok {
		r1 = rf(ctx, dogusToLoad)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockRemoteDoguRegistry_GetDogus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetDogus'
type MockRemoteDoguRegistry_GetDogus_Call struct {
	*mock.Call
}

// GetDogus is a helper method to define mock.On call
//   - ctx context.Context
//   - dogusToLoad []dogu.QualifiedVersion
func (_e *MockRemoteDoguRegistry_Expecter) GetDogus(ctx interface{}, dogusToLoad interface{}) *MockRemoteDoguRegistry_GetDogus_Call {
	return &MockRemoteDoguRegistry_GetDogus_Call{Call: _e.mock.On("GetDogus", ctx, dogusToLoad)}
}

func (_c *MockRemoteDoguRegistry_GetDogus_Call) Run(run func(ctx context.Context, dogusToLoad []dogu.QualifiedVersion)) *MockRemoteDoguRegistry_GetDogus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]dogu.QualifiedVersion))
	})
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogus_Call) Return(_a0 map[dogu.QualifiedName]*core.Dogu, _a1 error) *MockRemoteDoguRegistry_GetDogus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockRemoteDoguRegistry_GetDogus_Call) RunAndReturn(run func(context.Context, []dogu.QualifiedVersion) (map[dogu.QualifiedName]*core.Dogu, error)) *MockRemoteDoguRegistry_GetDogus_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockRemoteDoguRegistry creates a new instance of MockRemoteDoguRegistry. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRemoteDoguRegistry(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRemoteDoguRegistry {
	mock := &MockRemoteDoguRegistry{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
