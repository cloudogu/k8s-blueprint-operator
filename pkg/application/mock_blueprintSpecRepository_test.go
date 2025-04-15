// Code generated by mockery v2.53.3. DO NOT EDIT.

package application

import (
	context "context"

	domain "github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	mock "github.com/stretchr/testify/mock"
)

// mockBlueprintSpecRepository is an autogenerated mock type for the blueprintSpecRepository type
type mockBlueprintSpecRepository struct {
	mock.Mock
}

type mockBlueprintSpecRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *mockBlueprintSpecRepository) EXPECT() *mockBlueprintSpecRepository_Expecter {
	return &mockBlueprintSpecRepository_Expecter{mock: &_m.Mock}
}

// GetById provides a mock function with given fields: ctx, blueprintId
func (_m *mockBlueprintSpecRepository) GetById(ctx context.Context, blueprintId string) (*domain.BlueprintSpec, error) {
	ret := _m.Called(ctx, blueprintId)

	if len(ret) == 0 {
		panic("no return value specified for GetById")
	}

	var r0 *domain.BlueprintSpec
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.BlueprintSpec, error)); ok {
		return rf(ctx, blueprintId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.BlueprintSpec); ok {
		r0 = rf(ctx, blueprintId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.BlueprintSpec)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, blueprintId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockBlueprintSpecRepository_GetById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetById'
type mockBlueprintSpecRepository_GetById_Call struct {
	*mock.Call
}

// GetById is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintId string
func (_e *mockBlueprintSpecRepository_Expecter) GetById(ctx interface{}, blueprintId interface{}) *mockBlueprintSpecRepository_GetById_Call {
	return &mockBlueprintSpecRepository_GetById_Call{Call: _e.mock.On("GetById", ctx, blueprintId)}
}

func (_c *mockBlueprintSpecRepository_GetById_Call) Run(run func(ctx context.Context, blueprintId string)) *mockBlueprintSpecRepository_GetById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockBlueprintSpecRepository_GetById_Call) Return(_a0 *domain.BlueprintSpec, _a1 error) *mockBlueprintSpecRepository_GetById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockBlueprintSpecRepository_GetById_Call) RunAndReturn(run func(context.Context, string) (*domain.BlueprintSpec, error)) *mockBlueprintSpecRepository_GetById_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, blueprintSpec
func (_m *mockBlueprintSpecRepository) Update(ctx context.Context, blueprintSpec *domain.BlueprintSpec) error {
	ret := _m.Called(ctx, blueprintSpec)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.BlueprintSpec) error); ok {
		r0 = rf(ctx, blueprintSpec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockBlueprintSpecRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type mockBlueprintSpecRepository_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - blueprintSpec *domain.BlueprintSpec
func (_e *mockBlueprintSpecRepository_Expecter) Update(ctx interface{}, blueprintSpec interface{}) *mockBlueprintSpecRepository_Update_Call {
	return &mockBlueprintSpecRepository_Update_Call{Call: _e.mock.On("Update", ctx, blueprintSpec)}
}

func (_c *mockBlueprintSpecRepository_Update_Call) Run(run func(ctx context.Context, blueprintSpec *domain.BlueprintSpec)) *mockBlueprintSpecRepository_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*domain.BlueprintSpec))
	})
	return _c
}

func (_c *mockBlueprintSpecRepository_Update_Call) Return(_a0 error) *mockBlueprintSpecRepository_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockBlueprintSpecRepository_Update_Call) RunAndReturn(run func(context.Context, *domain.BlueprintSpec) error) *mockBlueprintSpecRepository_Update_Call {
	_c.Call.Return(run)
	return _c
}

// newMockBlueprintSpecRepository creates a new instance of mockBlueprintSpecRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockBlueprintSpecRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockBlueprintSpecRepository {
	mock := &mockBlueprintSpecRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
