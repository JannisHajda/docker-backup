package worker

import (
	"docker-backup/interfaces"
	"docker-backup/internal/borgclient"
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"docker-backup/internal/dockerclient"
	"fmt"
	"sync"
)

type Worker struct {
	db              interfaces.DatabaseClient
	dc              interfaces.DockerClient
	bc              interfaces.BorgClient
	workerContainer interfaces.DockerContainer
	sourceContainer interfaces.DockerContainer
	sourceVolumes   []interfaces.DockerVolume
	outputVolume    interfaces.DockerVolume
}

const (
	passphrase = "test"
	sourcePath = "/source"
	outputPath = "/output"
)

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

func (w *Worker) mountSourceVolumes(volumes []interfaces.DockerVolume) {
	for _, volume := range volumes {
		volume.SetMountPoint(fmt.Sprintf("%s/%s", sourcePath, volume.GetName()))
		w.sourceVolumes = append(w.sourceVolumes, volume)
	}
}

func (w *Worker) mountOutputVolume(volume interfaces.DockerVolume) {
	volume.SetMountPoint(outputPath)
	w.outputVolume = volume
}

func NewWorker(containerId string, outputVolumeName string) (interfaces.Worker, error) {
	dc := getDockerClient()

	source, err := dc.GetContainer(containerId)
	if err != nil {
		return nil, err
	}

	output, err := dc.CreateVolume(outputVolumeName)
	if err != nil {
		return nil, err
	}

	w := &Worker{
		dc:              dc,
		sourceContainer: source,
		outputVolume:    output,
	}

	sourceVolumes := w.getSourceVolumes(source)
	w.mountSourceVolumes(sourceVolumes)
	w.mountOutputVolume(output)

	w.workerContainer, err = w.dc.CreateContainer("worker", w.sourceVolumes, nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) Backup() error {
	if w.workerContainer == nil {
		return fmt.Errorf("worker container not created")
	}

	// pre-backup

	err := w.sourceContainer.Stop()
	if err != nil {
		return err
	}

	defer w.sourceContainer.Start()

	err = w.workerContainer.Start()
	if err != nil {
		return err
	}

	bc, err := borgclient.NewBorgClient(w.workerContainer, sourcePath, outputPath)
	if err != nil {
		return err
	}

	w.bc = bc

	defer w.workerContainer.StopAndRemove()

	errs := w.backupVolumes(w.sourceVolumes)
	if len(errs) > 0 {
		panic(errs[0])
	}

	// post-backup

	return nil
}

func (w *Worker) backupVolumes(volumes []interfaces.DockerVolume) []error {
	var errs []error
	var errsChan = make(chan error, len(volumes))
	var wg sync.WaitGroup

	for _, v := range volumes {
		wg.Add(1)

		go func(v interfaces.DockerVolume) {
			defer wg.Done()

			err := w.backupVolume(v)
			if err != nil {
				errsChan <- err
			}
		}(v)
	}

	// close errsChan when all volumes are backed u volumes are backed upp
	go func() {
		wg.Wait()
		close(errsChan)
	}()

	// collect errors
	for err := range errsChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (w *Worker) backupVolume(v interfaces.DockerVolume) error {
	repo, err := w.bc.GetOrCreateRepository(v.GetName(), passphrase)
	if err != nil {
		return err
	}

	return repo.Backup()
}

func (w *Worker) getSourceVolumes(c interfaces.DockerContainer) []interfaces.DockerVolume {
	volumes := c.GetVolumes()
	for _, volume := range volumes {
		path := fmt.Sprintf("%s/%s", sourcePath, volume.GetName())
		volume.SetMountPoint(path)
	}

	return volumes
}
