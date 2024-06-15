// Code generated by mockery v2.43.1. DO NOT EDIT.

package organization

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

type MockService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockService) EXPECT() *MockService_Expecter {
	return &MockService_Expecter{mock: &_m.Mock}
}

// CreateOrganization provides a mock function with given fields: ctx, subdomain, name
func (_m *MockService) CreateOrganization(ctx context.Context, subdomain string, name string) (Organization, error) {
	ret := _m.Called(ctx, subdomain, name)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrganization")
	}

	var r0 Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (Organization, error)); ok {
		return rf(ctx, subdomain, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) Organization); ok {
		r0 = rf(ctx, subdomain, name)
	} else {
		r0 = ret.Get(0).(Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, subdomain, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_CreateOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOrganization'
type MockService_CreateOrganization_Call struct {
	*mock.Call
}

// CreateOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - subdomain string
//   - name string
func (_e *MockService_Expecter) CreateOrganization(ctx interface{}, subdomain interface{}, name interface{}) *MockService_CreateOrganization_Call {
	return &MockService_CreateOrganization_Call{Call: _e.mock.On("CreateOrganization", ctx, subdomain, name)}
}

func (_c *MockService_CreateOrganization_Call) Run(run func(ctx context.Context, subdomain string, name string)) *MockService_CreateOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockService_CreateOrganization_Call) Return(_a0 Organization, _a1 error) *MockService_CreateOrganization_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_CreateOrganization_Call) RunAndReturn(run func(context.Context, string, string) (Organization, error)) *MockService_CreateOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOrganization provides a mock function with given fields: ctx, id, comment
func (_m *MockService) DeleteOrganization(ctx context.Context, id int64, comment string) error {
	ret := _m.Called(ctx, id, comment)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_DeleteOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOrganization'
type MockService_DeleteOrganization_Call struct {
	*mock.Call
}

// DeleteOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - comment string
func (_e *MockService_Expecter) DeleteOrganization(ctx interface{}, id interface{}, comment interface{}) *MockService_DeleteOrganization_Call {
	return &MockService_DeleteOrganization_Call{Call: _e.mock.On("DeleteOrganization", ctx, id, comment)}
}

func (_c *MockService_DeleteOrganization_Call) Run(run func(ctx context.Context, id int64, comment string)) *MockService_DeleteOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockService_DeleteOrganization_Call) Return(_a0 error) *MockService_DeleteOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_DeleteOrganization_Call) RunAndReturn(run func(context.Context, int64, string) error) *MockService_DeleteOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// DisableOrganization provides a mock function with given fields: ctx, id, comment
func (_m *MockService) DisableOrganization(ctx context.Context, id int64, comment string) error {
	ret := _m.Called(ctx, id, comment)

	if len(ret) == 0 {
		panic("no return value specified for DisableOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_DisableOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DisableOrganization'
type MockService_DisableOrganization_Call struct {
	*mock.Call
}

// DisableOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - comment string
func (_e *MockService_Expecter) DisableOrganization(ctx interface{}, id interface{}, comment interface{}) *MockService_DisableOrganization_Call {
	return &MockService_DisableOrganization_Call{Call: _e.mock.On("DisableOrganization", ctx, id, comment)}
}

func (_c *MockService_DisableOrganization_Call) Run(run func(ctx context.Context, id int64, comment string)) *MockService_DisableOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockService_DisableOrganization_Call) Return(_a0 error) *MockService_DisableOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_DisableOrganization_Call) RunAndReturn(run func(context.Context, int64, string) error) *MockService_DisableOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// EnableOrganization provides a mock function with given fields: ctx, id, comment
func (_m *MockService) EnableOrganization(ctx context.Context, id int64, comment string) error {
	ret := _m.Called(ctx, id, comment)

	if len(ret) == 0 {
		panic("no return value specified for EnableOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_EnableOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EnableOrganization'
type MockService_EnableOrganization_Call struct {
	*mock.Call
}

// EnableOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - comment string
func (_e *MockService_Expecter) EnableOrganization(ctx interface{}, id interface{}, comment interface{}) *MockService_EnableOrganization_Call {
	return &MockService_EnableOrganization_Call{Call: _e.mock.On("EnableOrganization", ctx, id, comment)}
}

func (_c *MockService_EnableOrganization_Call) Run(run func(ctx context.Context, id int64, comment string)) *MockService_EnableOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockService_EnableOrganization_Call) Return(_a0 error) *MockService_EnableOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_EnableOrganization_Call) RunAndReturn(run func(context.Context, int64, string) error) *MockService_EnableOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationByID provides a mock function with given fields: ctx, id
func (_m *MockService) GetOrganizationByID(ctx context.Context, id int64) (Organization, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationByID")
	}

	var r0 Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (Organization, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) Organization); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_GetOrganizationByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationByID'
type MockService_GetOrganizationByID_Call struct {
	*mock.Call
}

// GetOrganizationByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockService_Expecter) GetOrganizationByID(ctx interface{}, id interface{}) *MockService_GetOrganizationByID_Call {
	return &MockService_GetOrganizationByID_Call{Call: _e.mock.On("GetOrganizationByID", ctx, id)}
}

func (_c *MockService_GetOrganizationByID_Call) Run(run func(ctx context.Context, id int64)) *MockService_GetOrganizationByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockService_GetOrganizationByID_Call) Return(_a0 Organization, _a1 error) *MockService_GetOrganizationByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_GetOrganizationByID_Call) RunAndReturn(run func(context.Context, int64) (Organization, error)) *MockService_GetOrganizationByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationByName provides a mock function with given fields: ctx, name
func (_m *MockService) GetOrganizationByName(ctx context.Context, name string) (Organization, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationByName")
	}

	var r0 Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (Organization, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) Organization); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_GetOrganizationByName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationByName'
type MockService_GetOrganizationByName_Call struct {
	*mock.Call
}

// GetOrganizationByName is a helper method to define mock.On call
//   - ctx context.Context
//   - name string
func (_e *MockService_Expecter) GetOrganizationByName(ctx interface{}, name interface{}) *MockService_GetOrganizationByName_Call {
	return &MockService_GetOrganizationByName_Call{Call: _e.mock.On("GetOrganizationByName", ctx, name)}
}

func (_c *MockService_GetOrganizationByName_Call) Run(run func(ctx context.Context, name string)) *MockService_GetOrganizationByName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockService_GetOrganizationByName_Call) Return(_a0 Organization, _a1 error) *MockService_GetOrganizationByName_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_GetOrganizationByName_Call) RunAndReturn(run func(context.Context, string) (Organization, error)) *MockService_GetOrganizationByName_Call {
	_c.Call.Return(run)
	return _c
}

// GetOrganizationBySubdomain provides a mock function with given fields: ctx, subdomain
func (_m *MockService) GetOrganizationBySubdomain(ctx context.Context, subdomain string) (Organization, error) {
	ret := _m.Called(ctx, subdomain)

	if len(ret) == 0 {
		panic("no return value specified for GetOrganizationBySubdomain")
	}

	var r0 Organization
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (Organization, error)); ok {
		return rf(ctx, subdomain)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) Organization); ok {
		r0 = rf(ctx, subdomain)
	} else {
		r0 = ret.Get(0).(Organization)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, subdomain)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockService_GetOrganizationBySubdomain_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOrganizationBySubdomain'
type MockService_GetOrganizationBySubdomain_Call struct {
	*mock.Call
}

// GetOrganizationBySubdomain is a helper method to define mock.On call
//   - ctx context.Context
//   - subdomain string
func (_e *MockService_Expecter) GetOrganizationBySubdomain(ctx interface{}, subdomain interface{}) *MockService_GetOrganizationBySubdomain_Call {
	return &MockService_GetOrganizationBySubdomain_Call{Call: _e.mock.On("GetOrganizationBySubdomain", ctx, subdomain)}
}

func (_c *MockService_GetOrganizationBySubdomain_Call) Run(run func(ctx context.Context, subdomain string)) *MockService_GetOrganizationBySubdomain_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockService_GetOrganizationBySubdomain_Call) Return(_a0 Organization, _a1 error) *MockService_GetOrganizationBySubdomain_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockService_GetOrganizationBySubdomain_Call) RunAndReturn(run func(context.Context, string) (Organization, error)) *MockService_GetOrganizationBySubdomain_Call {
	_c.Call.Return(run)
	return _c
}

// SuspendOrganization provides a mock function with given fields: ctx, id, comment
func (_m *MockService) SuspendOrganization(ctx context.Context, id int64, comment string) error {
	ret := _m.Called(ctx, id, comment)

	if len(ret) == 0 {
		panic("no return value specified for SuspendOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_SuspendOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SuspendOrganization'
type MockService_SuspendOrganization_Call struct {
	*mock.Call
}

// SuspendOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - comment string
func (_e *MockService_Expecter) SuspendOrganization(ctx interface{}, id interface{}, comment interface{}) *MockService_SuspendOrganization_Call {
	return &MockService_SuspendOrganization_Call{Call: _e.mock.On("SuspendOrganization", ctx, id, comment)}
}

func (_c *MockService_SuspendOrganization_Call) Run(run func(ctx context.Context, id int64, comment string)) *MockService_SuspendOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockService_SuspendOrganization_Call) Return(_a0 error) *MockService_SuspendOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_SuspendOrganization_Call) RunAndReturn(run func(context.Context, int64, string) error) *MockService_SuspendOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// UnsuspendOrganization provides a mock function with given fields: ctx, id, comment
func (_m *MockService) UnsuspendOrganization(ctx context.Context, id int64, comment string) error {
	ret := _m.Called(ctx, id, comment)

	if len(ret) == 0 {
		panic("no return value specified for UnsuspendOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, comment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_UnsuspendOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnsuspendOrganization'
type MockService_UnsuspendOrganization_Call struct {
	*mock.Call
}

// UnsuspendOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - comment string
func (_e *MockService_Expecter) UnsuspendOrganization(ctx interface{}, id interface{}, comment interface{}) *MockService_UnsuspendOrganization_Call {
	return &MockService_UnsuspendOrganization_Call{Call: _e.mock.On("UnsuspendOrganization", ctx, id, comment)}
}

func (_c *MockService_UnsuspendOrganization_Call) Run(run func(ctx context.Context, id int64, comment string)) *MockService_UnsuspendOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockService_UnsuspendOrganization_Call) Return(_a0 error) *MockService_UnsuspendOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_UnsuspendOrganization_Call) RunAndReturn(run func(context.Context, int64, string) error) *MockService_UnsuspendOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOrganization provides a mock function with given fields: ctx, id, name
func (_m *MockService) UpdateOrganization(ctx context.Context, id int64, name string) error {
	ret := _m.Called(ctx, id, name)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrganization")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockService_UpdateOrganization_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOrganization'
type MockService_UpdateOrganization_Call struct {
	*mock.Call
}

// UpdateOrganization is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
//   - name string
func (_e *MockService_Expecter) UpdateOrganization(ctx interface{}, id interface{}, name interface{}) *MockService_UpdateOrganization_Call {
	return &MockService_UpdateOrganization_Call{Call: _e.mock.On("UpdateOrganization", ctx, id, name)}
}

func (_c *MockService_UpdateOrganization_Call) Run(run func(ctx context.Context, id int64, name string)) *MockService_UpdateOrganization_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(string))
	})
	return _c
}

func (_c *MockService_UpdateOrganization_Call) Return(_a0 error) *MockService_UpdateOrganization_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockService_UpdateOrganization_Call) RunAndReturn(run func(context.Context, int64, string) error) *MockService_UpdateOrganization_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockService creates a new instance of MockService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockService {
	mock := &MockService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
