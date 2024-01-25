package interfaces

type Table interface {
	Init() error
}

type ProjectsTable interface {
	Table
	GetByName(project_name string) (DatabaseProject, error)
	GetByID(project_id int64) (DatabaseProject, error)
	GetAll() ([]DatabaseProject, error)
	Add(p DatabaseProject) (DatabaseProject, error)
	Remove(p DatabaseProject) error
}

type ContainersTable interface {
	Table
	GetByName(container_name string) (DatabaseContainer, error)
	GetByID(container_id string) (DatabaseContainer, error)
	GetAll() ([]DatabaseContainer, error)
	Add(c DatabaseContainer) (DatabaseContainer, error)
	Remove(c DatabaseContainer) error
}

type ProjectContainersTable interface {
	Table
	GetByName(p DatabaseProject, container_name string) (DatabaseContainer, error)
	GetByID(p DatabaseProject, container_id string) (DatabaseContainer, error)
	GetAll(p DatabaseProject) ([]DatabaseContainer, error)
	Add(p DatabaseProject, c DatabaseContainer) error
	Remove(p DatabaseProject, c DatabaseContainer) error
}

type VolumesTable interface {
	Table
	GetByName(volume_name string) (DatabaseVolume, error)
	GetByID(volume_id int64) (DatabaseVolume, error)
	GetAll() ([]DatabaseVolume, error)
	Add(v DatabaseVolume) (DatabaseVolume, error)
	Remove(v DatabaseVolume) error
}

type ContainerVolumesTable interface {
	Table
	GetByName(c DatabaseContainer, volume_name string) (DatabaseVolume, error)
	GetByID(c DatabaseContainer, volume_id int64) (DatabaseVolume, error)
	GetAll(c DatabaseContainer) ([]DatabaseVolume, error)
	Add(c DatabaseContainer, v DatabaseVolume) error
	Remove(c DatabaseContainer, v DatabaseVolume) error
}
