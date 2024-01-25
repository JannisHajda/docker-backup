package db

import (
	"database/sql"
	"docker-backup/errors"
	"docker-backup/interfaces"

	_ "github.com/lib/pq"
)

type DatabaseClient struct {
	driver interfaces.DatabaseDriver
	pt     interfaces.ProjectsTable
	ct     interfaces.ContainersTable
	pct    interfaces.ProjectContainersTable
	vt     interfaces.VolumesTable
	cvt    interfaces.ContainerVolumesTable
}

func NewDatabaseClient(driver interfaces.DatabaseDriver) (interfaces.DatabaseClient, error) {
	client := &DatabaseClient{driver: driver}

	pt, err := NewProjectTable(client)
	if err != nil {
		return nil, err
	}

	ct, err := NewContainersTable(client)
	if err != nil {
		return nil, err
	}

	pct, err := NewProjectContainersTable(client)
	if err != nil {
		return nil, err
	}

	vt, err := NewVolumesTable(client)
	if err != nil {
		return nil, err
	}

	cvt, err := NewContainerVolumesTable(client)
	if err != nil {
		return nil, err
	}

	client.pt = pt
	client.ct = ct
	client.pct = pct
	client.vt = vt
	client.cvt = cvt

	return client, nil
}

func (d *DatabaseClient) GetDriver() interfaces.DatabaseDriver {
	return d.driver
}

func (d *DatabaseClient) SetDriver(driver interfaces.DatabaseDriver) {
	d.driver = driver
}

func (d *DatabaseClient) GetProjectsTable() interfaces.ProjectsTable {
	return d.pt
}

func (d *DatabaseClient) SetProjectsTable(pt interfaces.ProjectsTable) {
	d.pt = pt
}

func (d *DatabaseClient) GetContainersTable() interfaces.ContainersTable {
	return d.ct
}

func (d *DatabaseClient) SetContainersTable(ct interfaces.ContainersTable) {
	d.ct = ct
}

func (d *DatabaseClient) GetProjectContainersTable() interfaces.ProjectContainersTable {
	return d.pct
}

func (d *DatabaseClient) SetProjectContainersTable(pct interfaces.ProjectContainersTable) {
	d.pct = pct
}

func (d *DatabaseClient) GetVolumesTable() interfaces.VolumesTable {
	return d.vt
}

func (d *DatabaseClient) SetVolumesTable(vt interfaces.VolumesTable) {
	d.vt = vt
}

func (d *DatabaseClient) GetContainerVolumesTable() interfaces.ContainerVolumesTable {
	return d.cvt
}

func (d *DatabaseClient) SetContainerVolumesTable(cvt interfaces.ContainerVolumesTable) {
	d.cvt = cvt
}

func (d *DatabaseClient) AddProject(name string) (interfaces.DatabaseProject, error) {
	p := NewDatabaseProject(d, 0, name)
	p, err := d.pt.Add(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *DatabaseClient) RemoveProject(p interfaces.DatabaseProject) error {
	return d.pt.Remove(p)
}

func (d *DatabaseClient) GetProjects() ([]interfaces.DatabaseProject, error) {
	projects, err := d.pt.GetAll()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (d *DatabaseClient) GetProjectByID(id int64) (interfaces.DatabaseProject, error) {
	project, err := d.pt.GetByID(id)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (d *DatabaseClient) GetProjectByName(name string) (interfaces.DatabaseProject, error) {
	project, err := d.pt.GetByName(name)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (d *DatabaseClient) AddContainer(id string, name string) (interfaces.DatabaseContainer, error) {
	c := NewDatabaseContainer(d, id, name)
	c, err := d.ct.Add(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (d *DatabaseClient) RemoveContainer(c interfaces.DatabaseContainer) error {
	return d.ct.Remove(c)
}

func (d *DatabaseClient) GetContainers() ([]interfaces.DatabaseContainer, error) {
	containers, err := d.ct.GetAll()
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (d *DatabaseClient) GetOrAddContainer(id string, name string) (interfaces.DatabaseContainer, error) {
	container, err := d.ct.GetByName(name)

	if err != nil {
		if _, ok := err.(*errors.ContainerNotFoundError); ok {
			container, err = d.AddContainer(id, name)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return container, nil
}

func (d *DatabaseClient) GetContainerByID(id string) (interfaces.DatabaseContainer, error) {
	container, err := d.ct.GetByID(id)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (d *DatabaseClient) GetContainerByName(name string) (interfaces.DatabaseContainer, error) {
	container, err := d.ct.GetByName(name)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (d *DatabaseClient) AddVolume(name string) (interfaces.DatabaseVolume, error) {
	v := NewDatabaseVolume(d, 0, name)
	v, err := d.vt.Add(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (d *DatabaseClient) RemoveVolume(v interfaces.DatabaseVolume) error {
	return d.vt.Remove(v)
}

func (d *DatabaseClient) GetVolumes() ([]interfaces.DatabaseVolume, error) {
	volumes, err := d.vt.GetAll()
	if err != nil {
		return nil, err
	}

	return volumes, nil
}

func (d *DatabaseClient) GetVolumeByName(name string) (interfaces.DatabaseVolume, error) {
	volume, err := d.vt.GetByName(name)
	if err != nil {
		return nil, err
	}

	return volume, nil
}

func (d *DatabaseClient) GetVolumeByID(id int64) (interfaces.DatabaseVolume, error) {
	volume, err := d.vt.GetByID(id)
	if err != nil {
		return nil, err
	}

	return volume, nil
}

func (d *DatabaseClient) GetOrAddVolume(name string) (interfaces.DatabaseVolume, error) {
	volume, err := d.vt.GetByName(name)

	if err != nil {
		if _, ok := err.(*errors.VolumeNotFoundError); ok {
			volume, err = d.AddVolume(name)

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return volume, nil
}

func (d *DatabaseClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	res, err := d.driver.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DatabaseClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.driver.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
