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

	//	dbContainer, err := w.db.GetOrAddContainer(dockerContainer.GetID(), dockerContainer.GetName())
	//	if err != nil {
	//		return err
	//	}

	dockerVolumes := dockerContainer.GetVolumes()
	var errs []error
	for _, dockerVolume := range dockerVolumes {
		dockerVolume.SetMountPoint("/input/" + dockerVolume.GetName())

		//dbVolume, err := w.db.GetOrAddVolume(dockerVolume.GetName())
		//if err != nil {
		//errs = append(errs, err)
		//continue
		//}

		//err = dbContainer.AddVolume(dbVolume)
		//if err != nil {
		//errs = append(errs, err)
		//continue
		//}
	}

	dockerVolumes = append(dockerVolumes, w.output)

	if len(errs) > 0 {
		return errs[0]
	}

	_, err = dockerContainer.Exec("touch /test-volume1/test-file && touch /test-volume2/test-file")
	if err != nil {
		return err
	}

	_, err = dockerContainer.Exec("echo 'test' > /test-volume1/test-file")
	if err != nil {
		return err
	}

	//err = dockerContainer.Stop()
	//if err != nil {
	//	return err
	//}

	defer dockerContainer.Start()

	workerContainer, err := w.dc.CreateContainer("worker", dockerVolumes, nil)
	if err != nil {
		return err
	}

	err = workerContainer.Start()
	if err != nil {
		return err
	}

	defer workerContainer.StopAndRemove()

	bc, err := borgclient.NewBorgClient(workerContainer, "/input", "/output")
	if err != nil {
		return err
	}

	errs = []error{}
	for _, dockerVolume := range dockerVolumes {
		repo, err := bc.GetOrCreateRepository(dockerVolume.GetName(), "test")
		if err != nil {
			errs = append(errs, err)
			continue
		}

		err = repo.Backup()
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}

	return nil
}

func (w *Worker) BackupProject(projectName string) error {
	return nil
}
