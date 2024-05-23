// Code generated by mockery v2.42.2. DO NOT EDIT.

package auth

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ServiceMock is an autogenerated mock type for the Service type
type ServiceMock struct {
	mock.Mock
}

// Login provides a mock function with given fields: ctx, subdomain, email, password
func (_m *ServiceMock) Login(ctx context.Context, subdomain string, email string, password string) (string, error) {
	ret := _m.Called(ctx, subdomain, email, password)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (string, error)); ok {
		return rf(ctx, subdomain, email, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) string); ok {
		r0 = rf(ctx, subdomain, email, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, subdomain, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Register provides a mock function with given fields: ctx, email, password, subdomain, orgName
func (_m *ServiceMock) Register(ctx context.Context, email string, password string, subdomain string, orgName string) error {
	ret := _m.Called(ctx, email, password, subdomain, orgName)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) error); ok {
		r0 = rf(ctx, email, password, subdomain, orgName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewServiceMock creates a new instance of ServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ServiceMock {
	mock := &ServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
