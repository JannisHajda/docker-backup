// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DockerVolume is an autogenerated mock type for the DockerVolume type
type DockerVolume struct {
	mock.Mock
}

// GetMountPoint provides a mock function with given fields:
func (_m *DockerVolume) GetMountPoint() string {
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

// GetName provides a mock function with given fields:
func (_m *DockerVolume) GetName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetName")
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
func (_m *DockerVolume) IsRW() bool {
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

// SetMountPoint provides a mock function with given fields: mountPoint
func (_m *DockerVolume) SetMountPoint(mountPoint string) {
	_m.Called(mountPoint)
}

// NewDockerVolume creates a new instance of DockerVolume. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDockerVolume(t interface {
	mock.TestingT
	Cleanup(func())
}) *DockerVolume {
	mock := &DockerVolume{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}