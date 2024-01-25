package interfaces

import "database/sql"

type DatabaseClient interface {
	AddProject(name string) (DatabaseProject, error)
	RemoveProject(p DatabaseProject) error
	GetProjects() ([]DatabaseProject, error)
	GetProjectByID(id int64) (DatabaseProject, error)
	GetProjectByName(name string) (DatabaseProject, error)
	AddContainer(id string, name string) (DatabaseContainer, error)
	GetOrAddContainer(id string, name string) (DatabaseContainer, error)
	RemoveContainer(c DatabaseContainer) error
	GetContainers() ([]DatabaseContainer, error)
	GetContainerByID(id string) (DatabaseContainer, error)
	GetContainerByName(name string) (DatabaseContainer, error)
	AddVolume(name string) (DatabaseVolume, error)
	RemoveVolume(v DatabaseVolume) error
	GetVolumes() ([]DatabaseVolume, error)
	GetVolumeByID(id int64) (DatabaseVolume, error)
	GetVolumeByName(name string) (DatabaseVolume, error)
	GetOrAddVolume(name string) (DatabaseVolume, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	GetProjectsTable() ProjectsTable
	GetContainersTable() ContainersTable
	GetProjectContainersTable() ProjectContainersTable
	GetVolumesTable() VolumesTable
	GetContainerVolumesTable() ContainerVolumesTable
}

type DatabaseProject interface {
	GetID() int64
	GetName() string
	GetContainers() ([]DatabaseContainer, error)
	GetContainerByID(id string) (DatabaseContainer, error)
	GetContainerByName(name string) (DatabaseContainer, error)
	AddContainer(DatabaseContainer) error
	RemoveContainer(DatabaseContainer) error
}

type DatabaseContainer interface {
	GetID() string
	GetName() string
	GetVolumes() ([]DatabaseVolume, error)
	GetVolumeByID(id int64) (DatabaseVolume, error)
	GetVolumeByName(name string) (DatabaseVolume, error)
	AddVolume(DatabaseVolume) error
	RemoveVolume(DatabaseVolume) error
}

type DatabaseVolume interface {
	GetID() int64
	GetName() string
}

type DatabaseDriver interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
