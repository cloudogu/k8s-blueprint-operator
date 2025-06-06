// Code generated by mockery v2.53.3. DO NOT EDIT.

package restartcr

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	types "k8s.io/apimachinery/pkg/types"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"

	watch "k8s.io/apimachinery/pkg/watch"
)

// MockDoguRestartInterface is an autogenerated mock type for the DoguRestartInterface type
type MockDoguRestartInterface struct {
	mock.Mock
}

type MockDoguRestartInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDoguRestartInterface) EXPECT() *MockDoguRestartInterface_Expecter {
	return &MockDoguRestartInterface_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, dogu, opts
func (_m *MockDoguRestartInterface) Create(ctx context.Context, dogu *v2.DoguRestart, opts v1.CreateOptions) (*v2.DoguRestart, error) {
	ret := _m.Called(ctx, dogu, opts)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, v1.CreateOptions) (*v2.DoguRestart, error)); ok {
		return rf(ctx, dogu, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, v1.CreateOptions) *v2.DoguRestart); ok {
		r0 = rf(ctx, dogu, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v2.DoguRestart, v1.CreateOptions) error); ok {
		r1 = rf(ctx, dogu, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockDoguRestartInterface_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - dogu *v2.DoguRestart
//   - opts v1.CreateOptions
func (_e *MockDoguRestartInterface_Expecter) Create(ctx interface{}, dogu interface{}, opts interface{}) *MockDoguRestartInterface_Create_Call {
	return &MockDoguRestartInterface_Create_Call{Call: _e.mock.On("Create", ctx, dogu, opts)}
}

func (_c *MockDoguRestartInterface_Create_Call) Run(run func(ctx context.Context, dogu *v2.DoguRestart, opts v1.CreateOptions)) *MockDoguRestartInterface_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v2.DoguRestart), args[2].(v1.CreateOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_Create_Call) Return(_a0 *v2.DoguRestart, _a1 error) *MockDoguRestartInterface_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguRestartInterface_Create_Call) RunAndReturn(run func(context.Context, *v2.DoguRestart, v1.CreateOptions) (*v2.DoguRestart, error)) *MockDoguRestartInterface_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *MockDoguRestartInterface) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguRestartInterface_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockDoguRestartInterface_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts v1.DeleteOptions
func (_e *MockDoguRestartInterface_Expecter) Delete(ctx interface{}, name interface{}, opts interface{}) *MockDoguRestartInterface_Delete_Call {
	return &MockDoguRestartInterface_Delete_Call{Call: _e.mock.On("Delete", ctx, name, opts)}
}

func (_c *MockDoguRestartInterface_Delete_Call) Run(run func(ctx context.Context, name string, opts v1.DeleteOptions)) *MockDoguRestartInterface_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(v1.DeleteOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_Delete_Call) Return(_a0 error) *MockDoguRestartInterface_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguRestartInterface_Delete_Call) RunAndReturn(run func(context.Context, string, v1.DeleteOptions) error) *MockDoguRestartInterface_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCollection provides a mock function with given fields: ctx, opts, listOpts
func (_m *MockDoguRestartInterface) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	ret := _m.Called(ctx, opts, listOpts)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCollection")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.DeleteOptions, v1.ListOptions) error); ok {
		r0 = rf(ctx, opts, listOpts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDoguRestartInterface_DeleteCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCollection'
type MockDoguRestartInterface_DeleteCollection_Call struct {
	*mock.Call
}

// DeleteCollection is a helper method to define mock.On call
//   - ctx context.Context
//   - opts v1.DeleteOptions
//   - listOpts v1.ListOptions
func (_e *MockDoguRestartInterface_Expecter) DeleteCollection(ctx interface{}, opts interface{}, listOpts interface{}) *MockDoguRestartInterface_DeleteCollection_Call {
	return &MockDoguRestartInterface_DeleteCollection_Call{Call: _e.mock.On("DeleteCollection", ctx, opts, listOpts)}
}

func (_c *MockDoguRestartInterface_DeleteCollection_Call) Run(run func(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions)) *MockDoguRestartInterface_DeleteCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(v1.DeleteOptions), args[2].(v1.ListOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_DeleteCollection_Call) Return(_a0 error) *MockDoguRestartInterface_DeleteCollection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDoguRestartInterface_DeleteCollection_Call) RunAndReturn(run func(context.Context, v1.DeleteOptions, v1.ListOptions) error) *MockDoguRestartInterface_DeleteCollection_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, name, opts
func (_m *MockDoguRestartInterface) Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.DoguRestart, error) {
	ret := _m.Called(ctx, name, opts)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) (*v2.DoguRestart, error)); ok {
		return rf(ctx, name, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, v1.GetOptions) *v2.DoguRestart); ok {
		r0 = rf(ctx, name, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, v1.GetOptions) error); ok {
		r1 = rf(ctx, name, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockDoguRestartInterface_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - opts v1.GetOptions
func (_e *MockDoguRestartInterface_Expecter) Get(ctx interface{}, name interface{}, opts interface{}) *MockDoguRestartInterface_Get_Call {
	return &MockDoguRestartInterface_Get_Call{Call: _e.mock.On("Get", ctx, name, opts)}
}

func (_c *MockDoguRestartInterface_Get_Call) Run(run func(ctx context.Context, name string, opts v1.GetOptions)) *MockDoguRestartInterface_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(v1.GetOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_Get_Call) Return(_a0 *v2.DoguRestart, _a1 error) *MockDoguRestartInterface_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguRestartInterface_Get_Call) RunAndReturn(run func(context.Context, string, v1.GetOptions) (*v2.DoguRestart, error)) *MockDoguRestartInterface_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, opts
func (_m *MockDoguRestartInterface) List(ctx context.Context, opts v1.ListOptions) (*v2.DoguRestartList, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 *v2.DoguRestartList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.ListOptions) (*v2.DoguRestartList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1.ListOptions) *v2.DoguRestartList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestartList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockDoguRestartInterface_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - opts v1.ListOptions
func (_e *MockDoguRestartInterface_Expecter) List(ctx interface{}, opts interface{}) *MockDoguRestartInterface_List_Call {
	return &MockDoguRestartInterface_List_Call{Call: _e.mock.On("List", ctx, opts)}
}

func (_c *MockDoguRestartInterface_List_Call) Run(run func(ctx context.Context, opts v1.ListOptions)) *MockDoguRestartInterface_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(v1.ListOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_List_Call) Return(_a0 *v2.DoguRestartList, _a1 error) *MockDoguRestartInterface_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguRestartInterface_List_Call) RunAndReturn(run func(context.Context, v1.ListOptions) (*v2.DoguRestartList, error)) *MockDoguRestartInterface_List_Call {
	_c.Call.Return(run)
	return _c
}

// Patch provides a mock function with given fields: ctx, name, pt, data, opts, subresources
func (_m *MockDoguRestartInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (*v2.DoguRestart, error) {
	_va := make([]interface{}, len(subresources))
	for _i := range subresources {
		_va[_i] = subresources[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, name, pt, data, opts)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Patch")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, v1.PatchOptions, ...string) (*v2.DoguRestart, error)); ok {
		return rf(ctx, name, pt, data, opts, subresources...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.PatchType, []byte, v1.PatchOptions, ...string) *v2.DoguRestart); ok {
		r0 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.PatchType, []byte, v1.PatchOptions, ...string) error); ok {
		r1 = rf(ctx, name, pt, data, opts, subresources...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_Patch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Patch'
type MockDoguRestartInterface_Patch_Call struct {
	*mock.Call
}

// Patch is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
//   - pt types.PatchType
//   - data []byte
//   - opts v1.PatchOptions
//   - subresources ...string
func (_e *MockDoguRestartInterface_Expecter) Patch(ctx interface{}, name interface{}, pt interface{}, data interface{}, opts interface{}, subresources ...interface{}) *MockDoguRestartInterface_Patch_Call {
	return &MockDoguRestartInterface_Patch_Call{Call: _e.mock.On("Patch",
		append([]interface{}{ctx, name, pt, data, opts}, subresources...)...)}
}

func (_c *MockDoguRestartInterface_Patch_Call) Run(run func(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string)) *MockDoguRestartInterface_Patch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-5)
		for i, a := range args[5:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(types.PatchType), args[3].([]byte), args[4].(v1.PatchOptions), variadicArgs...)
	})
	return _c
}

func (_c *MockDoguRestartInterface_Patch_Call) Return(result *v2.DoguRestart, err error) *MockDoguRestartInterface_Patch_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *MockDoguRestartInterface_Patch_Call) RunAndReturn(run func(context.Context, string, types.PatchType, []byte, v1.PatchOptions, ...string) (*v2.DoguRestart, error)) *MockDoguRestartInterface_Patch_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, dogu, opts
func (_m *MockDoguRestartInterface) Update(ctx context.Context, dogu *v2.DoguRestart, opts v1.UpdateOptions) (*v2.DoguRestart, error) {
	ret := _m.Called(ctx, dogu, opts)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, v1.UpdateOptions) (*v2.DoguRestart, error)); ok {
		return rf(ctx, dogu, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, v1.UpdateOptions) *v2.DoguRestart); ok {
		r0 = rf(ctx, dogu, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v2.DoguRestart, v1.UpdateOptions) error); ok {
		r1 = rf(ctx, dogu, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockDoguRestartInterface_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - dogu *v2.DoguRestart
//   - opts v1.UpdateOptions
func (_e *MockDoguRestartInterface_Expecter) Update(ctx interface{}, dogu interface{}, opts interface{}) *MockDoguRestartInterface_Update_Call {
	return &MockDoguRestartInterface_Update_Call{Call: _e.mock.On("Update", ctx, dogu, opts)}
}

func (_c *MockDoguRestartInterface_Update_Call) Run(run func(ctx context.Context, dogu *v2.DoguRestart, opts v1.UpdateOptions)) *MockDoguRestartInterface_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v2.DoguRestart), args[2].(v1.UpdateOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_Update_Call) Return(_a0 *v2.DoguRestart, _a1 error) *MockDoguRestartInterface_Update_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguRestartInterface_Update_Call) RunAndReturn(run func(context.Context, *v2.DoguRestart, v1.UpdateOptions) (*v2.DoguRestart, error)) *MockDoguRestartInterface_Update_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSpecWithRetry provides a mock function with given fields: ctx, doguRestart, modifySpecFn, opts
func (_m *MockDoguRestartInterface) UpdateSpecWithRetry(ctx context.Context, doguRestart *v2.DoguRestart, modifySpecFn func(v2.DoguRestartSpec) v2.DoguRestartSpec, opts v1.UpdateOptions) (*v2.DoguRestart, error) {
	ret := _m.Called(ctx, doguRestart, modifySpecFn, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSpecWithRetry")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, func(v2.DoguRestartSpec) v2.DoguRestartSpec, v1.UpdateOptions) (*v2.DoguRestart, error)); ok {
		return rf(ctx, doguRestart, modifySpecFn, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, func(v2.DoguRestartSpec) v2.DoguRestartSpec, v1.UpdateOptions) *v2.DoguRestart); ok {
		r0 = rf(ctx, doguRestart, modifySpecFn, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v2.DoguRestart, func(v2.DoguRestartSpec) v2.DoguRestartSpec, v1.UpdateOptions) error); ok {
		r1 = rf(ctx, doguRestart, modifySpecFn, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_UpdateSpecWithRetry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSpecWithRetry'
type MockDoguRestartInterface_UpdateSpecWithRetry_Call struct {
	*mock.Call
}

// UpdateSpecWithRetry is a helper method to define mock.On call
//   - ctx context.Context
//   - doguRestart *v2.DoguRestart
//   - modifySpecFn func(v2.DoguRestartSpec) v2.DoguRestartSpec
//   - opts v1.UpdateOptions
func (_e *MockDoguRestartInterface_Expecter) UpdateSpecWithRetry(ctx interface{}, doguRestart interface{}, modifySpecFn interface{}, opts interface{}) *MockDoguRestartInterface_UpdateSpecWithRetry_Call {
	return &MockDoguRestartInterface_UpdateSpecWithRetry_Call{Call: _e.mock.On("UpdateSpecWithRetry", ctx, doguRestart, modifySpecFn, opts)}
}

func (_c *MockDoguRestartInterface_UpdateSpecWithRetry_Call) Run(run func(ctx context.Context, doguRestart *v2.DoguRestart, modifySpecFn func(v2.DoguRestartSpec) v2.DoguRestartSpec, opts v1.UpdateOptions)) *MockDoguRestartInterface_UpdateSpecWithRetry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v2.DoguRestart), args[2].(func(v2.DoguRestartSpec) v2.DoguRestartSpec), args[3].(v1.UpdateOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_UpdateSpecWithRetry_Call) Return(result *v2.DoguRestart, err error) *MockDoguRestartInterface_UpdateSpecWithRetry_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *MockDoguRestartInterface_UpdateSpecWithRetry_Call) RunAndReturn(run func(context.Context, *v2.DoguRestart, func(v2.DoguRestartSpec) v2.DoguRestartSpec, v1.UpdateOptions) (*v2.DoguRestart, error)) *MockDoguRestartInterface_UpdateSpecWithRetry_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatus provides a mock function with given fields: ctx, dogu, opts
func (_m *MockDoguRestartInterface) UpdateStatus(ctx context.Context, dogu *v2.DoguRestart, opts v1.UpdateOptions) (*v2.DoguRestart, error) {
	ret := _m.Called(ctx, dogu, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatus")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, v1.UpdateOptions) (*v2.DoguRestart, error)); ok {
		return rf(ctx, dogu, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, v1.UpdateOptions) *v2.DoguRestart); ok {
		r0 = rf(ctx, dogu, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v2.DoguRestart, v1.UpdateOptions) error); ok {
		r1 = rf(ctx, dogu, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_UpdateStatus_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatus'
type MockDoguRestartInterface_UpdateStatus_Call struct {
	*mock.Call
}

// UpdateStatus is a helper method to define mock.On call
//   - ctx context.Context
//   - dogu *v2.DoguRestart
//   - opts v1.UpdateOptions
func (_e *MockDoguRestartInterface_Expecter) UpdateStatus(ctx interface{}, dogu interface{}, opts interface{}) *MockDoguRestartInterface_UpdateStatus_Call {
	return &MockDoguRestartInterface_UpdateStatus_Call{Call: _e.mock.On("UpdateStatus", ctx, dogu, opts)}
}

func (_c *MockDoguRestartInterface_UpdateStatus_Call) Run(run func(ctx context.Context, dogu *v2.DoguRestart, opts v1.UpdateOptions)) *MockDoguRestartInterface_UpdateStatus_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v2.DoguRestart), args[2].(v1.UpdateOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_UpdateStatus_Call) Return(_a0 *v2.DoguRestart, _a1 error) *MockDoguRestartInterface_UpdateStatus_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguRestartInterface_UpdateStatus_Call) RunAndReturn(run func(context.Context, *v2.DoguRestart, v1.UpdateOptions) (*v2.DoguRestart, error)) *MockDoguRestartInterface_UpdateStatus_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateStatusWithRetry provides a mock function with given fields: ctx, doguRestart, modifyStatusFn, opts
func (_m *MockDoguRestartInterface) UpdateStatusWithRetry(ctx context.Context, doguRestart *v2.DoguRestart, modifyStatusFn func(v2.DoguRestartStatus) v2.DoguRestartStatus, opts v1.UpdateOptions) (*v2.DoguRestart, error) {
	ret := _m.Called(ctx, doguRestart, modifyStatusFn, opts)

	if len(ret) == 0 {
		panic("no return value specified for UpdateStatusWithRetry")
	}

	var r0 *v2.DoguRestart
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, func(v2.DoguRestartStatus) v2.DoguRestartStatus, v1.UpdateOptions) (*v2.DoguRestart, error)); ok {
		return rf(ctx, doguRestart, modifyStatusFn, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v2.DoguRestart, func(v2.DoguRestartStatus) v2.DoguRestartStatus, v1.UpdateOptions) *v2.DoguRestart); ok {
		r0 = rf(ctx, doguRestart, modifyStatusFn, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v2.DoguRestart)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v2.DoguRestart, func(v2.DoguRestartStatus) v2.DoguRestartStatus, v1.UpdateOptions) error); ok {
		r1 = rf(ctx, doguRestart, modifyStatusFn, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_UpdateStatusWithRetry_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateStatusWithRetry'
type MockDoguRestartInterface_UpdateStatusWithRetry_Call struct {
	*mock.Call
}

// UpdateStatusWithRetry is a helper method to define mock.On call
//   - ctx context.Context
//   - doguRestart *v2.DoguRestart
//   - modifyStatusFn func(v2.DoguRestartStatus) v2.DoguRestartStatus
//   - opts v1.UpdateOptions
func (_e *MockDoguRestartInterface_Expecter) UpdateStatusWithRetry(ctx interface{}, doguRestart interface{}, modifyStatusFn interface{}, opts interface{}) *MockDoguRestartInterface_UpdateStatusWithRetry_Call {
	return &MockDoguRestartInterface_UpdateStatusWithRetry_Call{Call: _e.mock.On("UpdateStatusWithRetry", ctx, doguRestart, modifyStatusFn, opts)}
}

func (_c *MockDoguRestartInterface_UpdateStatusWithRetry_Call) Run(run func(ctx context.Context, doguRestart *v2.DoguRestart, modifyStatusFn func(v2.DoguRestartStatus) v2.DoguRestartStatus, opts v1.UpdateOptions)) *MockDoguRestartInterface_UpdateStatusWithRetry_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v2.DoguRestart), args[2].(func(v2.DoguRestartStatus) v2.DoguRestartStatus), args[3].(v1.UpdateOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_UpdateStatusWithRetry_Call) Return(result *v2.DoguRestart, err error) *MockDoguRestartInterface_UpdateStatusWithRetry_Call {
	_c.Call.Return(result, err)
	return _c
}

func (_c *MockDoguRestartInterface_UpdateStatusWithRetry_Call) RunAndReturn(run func(context.Context, *v2.DoguRestart, func(v2.DoguRestartStatus) v2.DoguRestartStatus, v1.UpdateOptions) (*v2.DoguRestart, error)) *MockDoguRestartInterface_UpdateStatusWithRetry_Call {
	_c.Call.Return(run)
	return _c
}

// Watch provides a mock function with given fields: ctx, opts
func (_m *MockDoguRestartInterface) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, opts)

	if len(ret) == 0 {
		panic("no return value specified for Watch")
	}

	var r0 watch.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.ListOptions) (watch.Interface, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDoguRestartInterface_Watch_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Watch'
type MockDoguRestartInterface_Watch_Call struct {
	*mock.Call
}

// Watch is a helper method to define mock.On call
//   - ctx context.Context
//   - opts v1.ListOptions
func (_e *MockDoguRestartInterface_Expecter) Watch(ctx interface{}, opts interface{}) *MockDoguRestartInterface_Watch_Call {
	return &MockDoguRestartInterface_Watch_Call{Call: _e.mock.On("Watch", ctx, opts)}
}

func (_c *MockDoguRestartInterface_Watch_Call) Run(run func(ctx context.Context, opts v1.ListOptions)) *MockDoguRestartInterface_Watch_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(v1.ListOptions))
	})
	return _c
}

func (_c *MockDoguRestartInterface_Watch_Call) Return(_a0 watch.Interface, _a1 error) *MockDoguRestartInterface_Watch_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDoguRestartInterface_Watch_Call) RunAndReturn(run func(context.Context, v1.ListOptions) (watch.Interface, error)) *MockDoguRestartInterface_Watch_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDoguRestartInterface creates a new instance of MockDoguRestartInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDoguRestartInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDoguRestartInterface {
	mock := &MockDoguRestartInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
