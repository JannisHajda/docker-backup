// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DockerBind is an autogenerated mock type for the DockerBind type
type DockerBind struct {
	mock.Mock
}

// GetHostPath provides a mock function with given fields:
func (_m *DockerBind) GetHostPath() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetHostPath")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetMountPoint provides a mock function with given fields:
func (_m *DockerBind) GetMountPoint() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetMountPoint")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsRW provides a mock function with given fields:
func (_m *DockerBind) IsRW() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsRW")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// NewDockerBind creates a new instance of DockerBind. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDockerBind(t interface {
	mock.TestingT
	Cleanup(func())
}) *DockerBind {
	mock := &DockerBind{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}