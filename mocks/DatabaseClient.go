// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	interfaces "docker-backup/interfaces"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// DatabaseClient is an autogenerated mock type for the DatabaseClient type
type DatabaseClient struct {
	mock.Mock
}

// AddContainer provides a mock function with given fields: id, name
func (_m *DatabaseClient) AddContainer(id string, name string) (interfaces.DatabaseContainer, error) {
	ret := _m.Called(id, name)

	if len(ret) == 0 {
		panic("no return value specified for AddContainer")
	}

	var r0 interfaces.DatabaseContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (interfaces.DatabaseContainer, error)); ok {
		return rf(id, name)
	}
	if rf, ok := ret.Get(0).(func(string, string) interfaces.DatabaseContainer); ok {
		r0 = rf(id, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddProject provides a mock function with given fields: name
func (_m *DatabaseClient) AddProject(name string) (interfaces.DatabaseProject, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for AddProject")
	}

	var r0 interfaces.DatabaseProject
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseProject, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseProject); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseProject)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AddVolume provides a mock function with given fields: name
func (_m *DatabaseClient) AddVolume(name string) (interfaces.DatabaseVolume, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for AddVolume")
	}

	var r0 interfaces.DatabaseVolume
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseVolume, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseVolume); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseVolume)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Exec provides a mock function with given fields: query, args
func (_m *DatabaseClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (sql.Result, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) sql.Result); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContainerByID provides a mock function with given fields: id
func (_m *DatabaseClient) GetContainerByID(id string) (interfaces.DatabaseContainer, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetContainerByID")
	}

	var r0 interfaces.DatabaseContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseContainer, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseContainer); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContainerByName provides a mock function with given fields: name
func (_m *DatabaseClient) GetContainerByName(name string) (interfaces.DatabaseContainer, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetContainerByName")
	}

	var r0 interfaces.DatabaseContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseContainer, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseContainer); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContainerVolumesTable provides a mock function with given fields:
func (_m *DatabaseClient) GetContainerVolumesTable() interfaces.ContainerVolumesTable {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetContainerVolumesTable")
	}

	var r0 interfaces.ContainerVolumesTable
	if rf, ok := ret.Get(0).(func() interfaces.ContainerVolumesTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.ContainerVolumesTable)
		}
	}

	return r0
}

// GetContainers provides a mock function with given fields:
func (_m *DatabaseClient) GetContainers() ([]interfaces.DatabaseContainer, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetContainers")
	}

	var r0 []interfaces.DatabaseContainer
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]interfaces.DatabaseContainer, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []interfaces.DatabaseContainer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interfaces.DatabaseContainer)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetContainersTable provides a mock function with given fields:
func (_m *DatabaseClient) GetContainersTable() interfaces.ContainersTable {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetContainersTable")
	}

	var r0 interfaces.ContainersTable
	if rf, ok := ret.Get(0).(func() interfaces.ContainersTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.ContainersTable)
		}
	}

	return r0
}

// GetOrAddContainer provides a mock function with given fields: id, name
func (_m *DatabaseClient) GetOrAddContainer(id string, name string) (interfaces.DatabaseContainer, error) {
	ret := _m.Called(id, name)

	if len(ret) == 0 {
		panic("no return value specified for GetOrAddContainer")
	}

	var r0 interfaces.DatabaseContainer
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (interfaces.DatabaseContainer, error)); ok {
		return rf(id, name)
	}
	if rf, ok := ret.Get(0).(func(string, string) interfaces.DatabaseContainer); ok {
		r0 = rf(id, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseContainer)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrAddVolume provides a mock function with given fields: name
func (_m *DatabaseClient) GetOrAddVolume(name string) (interfaces.DatabaseVolume, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetOrAddVolume")
	}

	var r0 interfaces.DatabaseVolume
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseVolume, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseVolume); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseVolume)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProjectByID provides a mock function with given fields: id
func (_m *DatabaseClient) GetProjectByID(id int64) (interfaces.DatabaseProject, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetProjectByID")
	}

	var r0 interfaces.DatabaseProject
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (interfaces.DatabaseProject, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) interfaces.DatabaseProject); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseProject)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProjectByName provides a mock function with given fields: name
func (_m *DatabaseClient) GetProjectByName(name string) (interfaces.DatabaseProject, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetProjectByName")
	}

	var r0 interfaces.DatabaseProject
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseProject, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseProject); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseProject)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProjectContainersTable provides a mock function with given fields:
func (_m *DatabaseClient) GetProjectContainersTable() interfaces.ProjectContainersTable {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetProjectContainersTable")
	}

	var r0 interfaces.ProjectContainersTable
	if rf, ok := ret.Get(0).(func() interfaces.ProjectContainersTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.ProjectContainersTable)
		}
	}

	return r0
}

// GetProjects provides a mock function with given fields:
func (_m *DatabaseClient) GetProjects() ([]interfaces.DatabaseProject, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetProjects")
	}

	var r0 []interfaces.DatabaseProject
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]interfaces.DatabaseProject, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []interfaces.DatabaseProject); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interfaces.DatabaseProject)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProjectsTable provides a mock function with given fields:
func (_m *DatabaseClient) GetProjectsTable() interfaces.ProjectsTable {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetProjectsTable")
	}

	var r0 interfaces.ProjectsTable
	if rf, ok := ret.Get(0).(func() interfaces.ProjectsTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.ProjectsTable)
		}
	}

	return r0
}

// GetVolumeByID provides a mock function with given fields: id
func (_m *DatabaseClient) GetVolumeByID(id int64) (interfaces.DatabaseVolume, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetVolumeByID")
	}

	var r0 interfaces.DatabaseVolume
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (interfaces.DatabaseVolume, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) interfaces.DatabaseVolume); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseVolume)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVolumeByName provides a mock function with given fields: name
func (_m *DatabaseClient) GetVolumeByName(name string) (interfaces.DatabaseVolume, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for GetVolumeByName")
	}

	var r0 interfaces.DatabaseVolume
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (interfaces.DatabaseVolume, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) interfaces.DatabaseVolume); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.DatabaseVolume)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVolumes provides a mock function with given fields:
func (_m *DatabaseClient) GetVolumes() ([]interfaces.DatabaseVolume, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetVolumes")
	}

	var r0 []interfaces.DatabaseVolume
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]interfaces.DatabaseVolume, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []interfaces.DatabaseVolume); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interfaces.DatabaseVolume)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVolumesTable provides a mock function with given fields:
func (_m *DatabaseClient) GetVolumesTable() interfaces.VolumesTable {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetVolumesTable")
	}

	var r0 interfaces.VolumesTable
	if rf, ok := ret.Get(0).(func() interfaces.VolumesTable); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interfaces.VolumesTable)
		}
	}

	return r0
}

// Query provides a mock function with given fields: query, args
func (_m *DatabaseClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 *sql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (*sql.Rows, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveContainer provides a mock function with given fields: c
func (_m *DatabaseClient) RemoveContainer(c interfaces.DatabaseContainer) error {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for RemoveContainer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interfaces.DatabaseContainer) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveProject provides a mock function with given fields: p
func (_m *DatabaseClient) RemoveProject(p interfaces.DatabaseProject) error {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for RemoveProject")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interfaces.DatabaseProject) error); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveVolume provides a mock function with given fields: v
func (_m *DatabaseClient) RemoveVolume(v interfaces.DatabaseVolume) error {
	ret := _m.Called(v)

	if len(ret) == 0 {
		panic("no return value specified for RemoveVolume")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(interfaces.DatabaseVolume) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDatabaseClient creates a new instance of DatabaseClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDatabaseClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *DatabaseClient {
	mock := &DatabaseClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
