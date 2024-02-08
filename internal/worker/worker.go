package worker

import (
	"docker-backup/interfaces"
	"docker-backup/internal/borgclient"
	"docker-backup/internal/db"
	"docker-backup/internal/db/driver"
	"docker-backup/internal/dockerclient"
	"fmt"
	"regexp"
)

const (
	DOCKER_IMAGE = "worker"
	sourcePath   = "/source"
	outputPath   = "/output"
)

type Worker struct {
	db interfaces.DatabaseClient
	dc interfaces.DockerClient
	bc interfaces.BorgClient

	workerContainer interfaces.DockerContainer
	sourceContainer interfaces.DockerContainer

	inputVolumes  []interfaces.DockerVolume
	localBackups  []interfaces.LocalBackup
	remoteBackups []interfaces.RemoteBackup

	passphrase string
}

func isRemoteHostNotFoundError(output string) bool {
	re := regexp.MustCompile("Name or service not known")
	matches := re.FindStringSubmatch(output)

	if len(matches) > 0 {
		return true
	}

	return false
}

func (w *Worker) initSSH(hosts []string, keyfiles []string) error {
	_, err := w.workerContainer.Exec("mkdir -p ~/.ssh")
	if err != nil {
		return err
	}

	_, err = w.workerContainer.Exec("eval `ssh-agent`")
	if err != nil {
		return err
	}

	for _, host := range hosts {
		err = w.addKnownHost(host)
		if err != nil {
			return err
		}
	}

	for _, keyfile := range keyfiles {
		err = w.addPrivateKey(keyfile)
		if err != nil {
			return err
		}
	}

	return err
}

func (w *Worker) addKnownHost(host string) error {
	_, err := w.workerContainer.Exec(fmt.Sprintf("ssh-keyscan -H %s >> ~/.ssh/known_hosts", host))
	return err
}

func (w *Worker) addPrivateKey(keyfile string) error {
	_, err := w.workerContainer.Exec(fmt.Sprintf("cp %s ~/.ssh", keyfile))
	return err
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

func (w *Worker) mountSourceVolumes(volumes []interfaces.DockerVolume) []interfaces.DockerVolume {
	for _, volume := range volumes {
		path := fmt.Sprintf("%s/%s", sourcePath, volume.GetName())
		volume.SetMountPoint(path)
	}

	return volumes
}

func (w *Worker) mountOutputVolumes(volumes []interfaces.DockerVolume) []interfaces.DockerVolume {
	for _, volume := range volumes {
		path := fmt.Sprintf("%s/%s", outputPath, volume.GetName())
		volume.SetMountPoint(path)
	}

	return volumes
}

func NewWorker(containerId string, passphrase string, localBackups []interfaces.LocalBackup, remoteBackups []interfaces.RemoteBackup) (interfaces.Worker, error) {
	dc := getDockerClient()

	container, err := dc.GetContainer(containerId)
	if err != nil {
		return nil, err
	}

	var errs []error
	var outputVolumes []interfaces.DockerVolume
	var validLocalBackups []interfaces.LocalBackup

	for _, backup := range localBackups {
		volume, err := dc.CreateVolume(backup.GetVolumeName())

		if err != nil {
			errs = append(errs, err)

		} else {
			outputVolumes = append(outputVolumes, volume)
			validLocalBackups = append(validLocalBackups, backup)
		}
	}

	if len(errs) > 0 {
		fmt.Sprintf("Failed to create output volumes: %v", errs)
	}

	inputVolumes := getSourceVolumes(container)

	w := &Worker{
		dc:              dc,
		sourceContainer: container,
		inputVolumes:    inputVolumes,
		localBackups:    validLocalBackups,
		remoteBackups:   remoteBackups,
		passphrase:      passphrase,
	}

	outputVolumes = w.mountOutputVolumes(outputVolumes)
	inputVolumes = w.mountSourceVolumes(inputVolumes)

	w.workerContainer, err = w.dc.CreateContainer(DOCKER_IMAGE, append(inputVolumes, outputVolumes...), nil)
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

	//err := w.sourceContainer.Stop()
	//if err != nil {
	//	return err
	//}

	//defer w.sourceContainer.Start()

	err := w.workerContainer.Start()
	if err != nil {
		return err
	}

	bc, err := borgclient.NewBorgClient(w.workerContainer)
	if err != nil {
		return err
	}

	w.bc = bc

	for _, backupLocation := range w.localBackups {
		backupDir := fmt.Sprintf("%s/%s", outputPath, backupLocation.GetVolumeName())

		var errs []error
		for _, volume := range w.inputVolumes {
			backupRepo := fmt.Sprintf("%s/%s", backupDir, volume.GetName())
			repo, err := w.bc.GetOrCreateRepository(backupRepo, w.passphrase)

			if err != nil {
				errs = append(errs, err)
			} else {
				err = repo.Backup(volume.GetMountPoint())
				if err != nil {
					errs = append(errs, err)
				}
			}
		}

		if len(errs) > 0 {
			return fmt.Errorf("Failed to backup volumes: %v", errs)
		}
	}

	if len(w.remoteBackups) > 0 {
		err = w.initSSH([]string{"test"}, []string{"test"})
		if err != nil {
			return err
		}

		for _, backupLocation := range w.remoteBackups {
			backupDir := fmt.Sprintf("%s@%s:%s/%s", backupLocation.GetUser(), backupLocation.GetHost(), backupLocation.GetPath())

			var errs []error
			for _, volume := range w.inputVolumes {
				backupRepo := fmt.Sprintf("%s/%s", backupDir, volume.GetName())
				repo, err := w.bc.GetOrCreateRepository(backupRepo, w.passphrase)

				if err != nil {
					errs = append(errs, err)
				} else {
					err = repo.Backup(volume.GetMountPoint())
					if err != nil {
						errs = append(errs, err)
					}
				}
			}

			if len(errs) > 0 {
				return fmt.Errorf("Failed to backup volumes: %v", errs)
			}
		}
	}

	defer w.workerContainer.StopAndRemove()

	// post-backup

	return nil
}

//func (w *Worker) backupVolumes(volumes []interfaces.DockerVolume) []error {
//	var errs []error
//	var errsChan = make(chan error, len(volumes))
//	var wg sync.WaitGroup
//
//	for _, v := range volumes {
//		wg.Add(1)
//
//		go func(v interfaces.DockerVolume) {
//			defer wg.Done()
//
//			err := w.backupVolume(v)
//			if err != nil {
//				errsChan <- err
//			}
//		}(v)
//	}
//
//	// close errsChan when all volumes are backed u volumes are backed upp
//	go func() {
//		wg.Wait()
//		close(errsChan)
//	}()
//
//	// collect errors
//	for err := range errsChan {
//		if err != nil {
//			errs = append(errs, err)
//		}
//	}
//
//	return errs
//}
//
//func (w *Worker) backupVolume(v interfaces.DockerVolume) error {
//	repo, err := w.bc.GetOrCreateRepository(v.GetName(), w.passphrase)
//	if err != nil {
//		return err
//	}
//
//	return repo.Backup()
//}

func getSourceVolumes(c interfaces.DockerContainer) []interfaces.DockerVolume {
	volumes := c.GetVolumes()
	for _, volume := range volumes {
		path := fmt.Sprintf("%s/%s", sourcePath, volume.GetName())
		volume.SetMountPoint(path)
	}

	return volumes
}

func (w *Worker) Stop() error {
	if w.workerContainer != nil {
		return w.workerContainer.StopAndRemove()
	}

	return nil
}
