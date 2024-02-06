package worker

import (
	"docker-backup/interfaces"
	"docker-backup/internal/borgclient"
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"docker-backup/internal/dockerclient"
)

type Worker struct {
	db        interfaces.DatabaseClient
	dc        interfaces.DockerClient
	container interfaces.DockerContainer
	output    interfaces.DockerVolume
}

func getDbClient() interfaces.DatabaseClient {
	driver, err := driver.NewPostgresDriver("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	client, err := db.NewDatabaseClient(driver)
	if err != nil {
		panic(err)
	}

	return client
}

func getDockerClient() interfaces.DockerClient {
	client, err := dockerclient.NewDockerClient()
	if err != nil {
		panic(err)
	}

	return client
}

func NewWorker(output string) (interfaces.Worker, error) {
	dockerClient := getDockerClient()
	outputVolume, err := dockerClient.CreateVolume(output)
	outputVolume.SetMountPoint("/output")

	if err != nil {
		return nil, err
	}

	return &Worker{
		dc:     dockerClient,
		output: outputVolume,
	}, nil
}

func (w *Worker) BackupContainer(containerIdentifier string) error {
	dockerContainer, err := w.dc.GetContainer(containerIdentifier)
	if err != nil {
		return err
	}

	// pre-backup

	dockerVolumes := w.getSourceVolumes(dockerContainer)

	err = dockerContainer.Stop()
	if err != nil {
		return err
	}

	defer dockerContainer.Start()

	workerContainer, err := w.createAndStartWorkerContainer(dockerVolumes, []interfaces.DockerBind{})
	defer workerContainer.StopAndRemove()

	bc, err := borgclient.NewBorgClient(workerContainer, "/input", "/output")
	if err != nil {
		return err
	}

	errs := w.backupVolumes(bc, dockerVolumes)
	if len(errs) > 0 {
		panic(errs[0])
	}

	// post-backup

	return nil
}

func (w *Worker) backupVolumes(bc interfaces.BorgClient, v []interfaces.DockerVolume) []error {
	var errs []error
	for _, volume := range v {
		err := w.backupVolume(bc, volume)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (w *Worker) backupVolume(bc interfaces.BorgClient, v interfaces.DockerVolume) error {
	repo, err := bc.GetOrCreateRepository(v.GetName(), "test")
	if err != nil {
		return err
	}

	return repo.Backup()
}

func (w *Worker) getSourceVolumes(c interfaces.DockerContainer) []interfaces.DockerVolume {
	volumes := c.GetVolumes()
	for _, volume := range volumes {
		volume.SetMountPoint("/input/" + volume.GetName())
	}

	return volumes
}

func (w *Worker) createAndStartWorkerContainer(volumes []interfaces.DockerVolume, binds []interfaces.DockerBind) (interfaces.DockerContainer, error) {
	c, err := w.dc.CreateContainer("worker", volumes, binds)
	if err != nil {
		return nil, err
	}

	err = c.Start()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (w *Worker) BackupProject(projectName string) error {
	return nil
}
