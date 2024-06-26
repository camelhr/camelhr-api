// Code generated by mockery v2.43.1. DO NOT EDIT.

package database

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

type MockDatabase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabase) EXPECT() *MockDatabase_Expecter {
	return &MockDatabase_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, dest, query, args
func (_m *MockDatabase) Exec(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockDatabase_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - dest interface{}
//   - query string
//   - args ...interface{}
func (_e *MockDatabase_Expecter) Exec(ctx interface{}, dest interface{}, query interface{}, args ...interface{}) *MockDatabase_Exec_Call {
	return &MockDatabase_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{ctx, dest, query}, args...)...)}
}

func (_c *MockDatabase_Exec_Call) Run(run func(ctx context.Context, dest interface{}, query string, args ...interface{})) *MockDatabase_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDatabase_Exec_Call) Return(_a0 error) *MockDatabase_Exec_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Exec_Call) RunAndReturn(run func(context.Context, interface{}, string, ...interface{}) error) *MockDatabase_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, dest, query, args
func (_m *MockDatabase) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockDatabase_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - dest interface{}
//   - query string
//   - args ...interface{}
func (_e *MockDatabase_Expecter) Get(ctx interface{}, dest interface{}, query interface{}, args ...interface{}) *MockDatabase_Get_Call {
	return &MockDatabase_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{ctx, dest, query}, args...)...)}
}

func (_c *MockDatabase_Get_Call) Run(run func(ctx context.Context, dest interface{}, query string, args ...interface{})) *MockDatabase_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDatabase_Get_Call) Return(_a0 error) *MockDatabase_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_Get_Call) RunAndReturn(run func(context.Context, interface{}, string, ...interface{}) error) *MockDatabase_Get_Call {
	_c.Call.Return(run)
	return _c
}

// List provides a mock function with given fields: ctx, dest, query, args
func (_m *MockDatabase) List(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type MockDatabase_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//   - ctx context.Context
//   - dest interface{}
//   - query string
//   - args ...interface{}
func (_e *MockDatabase_Expecter) List(ctx interface{}, dest interface{}, query interface{}, args ...interface{}) *MockDatabase_List_Call {
	return &MockDatabase_List_Call{Call: _e.mock.On("List",
		append([]interface{}{ctx, dest, query}, args...)...)}
}

func (_c *MockDatabase_List_Call) Run(run func(ctx context.Context, dest interface{}, query string, args ...interface{})) *MockDatabase_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDatabase_List_Call) Return(_a0 error) *MockDatabase_List_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_List_Call) RunAndReturn(run func(context.Context, interface{}, string, ...interface{}) error) *MockDatabase_List_Call {
	_c.Call.Return(run)
	return _c
}

// WithTx provides a mock function with given fields: ctx, txFn
func (_m *MockDatabase) WithTx(ctx context.Context, txFn func(context.Context) error) error {
	ret := _m.Called(ctx, txFn)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, txFn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_WithTx_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithTx'
type MockDatabase_WithTx_Call struct {
	*mock.Call
}

// WithTx is a helper method to define mock.On call
//   - ctx context.Context
//   - txFn func(context.Context) error
func (_e *MockDatabase_Expecter) WithTx(ctx interface{}, txFn interface{}) *MockDatabase_WithTx_Call {
	return &MockDatabase_WithTx_Call{Call: _e.mock.On("WithTx", ctx, txFn)}
}

func (_c *MockDatabase_WithTx_Call) Run(run func(ctx context.Context, txFn func(context.Context) error)) *MockDatabase_WithTx_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *MockDatabase_WithTx_Call) Return(_a0 error) *MockDatabase_WithTx_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_WithTx_Call) RunAndReturn(run func(context.Context, func(context.Context) error) error) *MockDatabase_WithTx_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
