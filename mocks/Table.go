// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Table is an autogenerated mock type for the Table type
type Table struct {
	mock.Mock
}

// Init provides a mock function with given fields:
func (_m *Table) Init() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTable creates a new instance of Table. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTable(t interface {
	mock.TestingT
	Cleanup(func())
}) *Table {
	mock := &Table{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}