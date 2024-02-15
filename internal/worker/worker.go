package worker

import (
	"docker-backup/interfaces"
	"docker-backup/internal/borg"
	"docker-backup/internal/docker"
	"docker-backup/internal/helper"
	"docker-backup/internal/ssh"
	"fmt"
	"strings"
	"sync"
)

const (
	DOCKER_IMAGE = "worker"
	inputPath    = "/input"
	outputPath   = "/output"
	keyfilesPath = "/keyfiles"
)

type Worker struct {
	db  interfaces.DatabaseClient
	dc  interfaces.DockerClient
	bc  interfaces.BorgClient
	ssh interfaces.SSHClient

	workerContainer interfaces.DockerContainer
	sourceContainer interfaces.DockerContainer

	inputVolumes  []interfaces.DockerVolume
	outputVolumes []interfaces.DockerVolume
	keyfiles      []interfaces.DockerBind
	localBackups  []interfaces.LocalBackup
	remoteBackups []interfaces.RemoteBackup

	passphrase string
}

func NewWorker(containerId string, passphrase string, localBackups []interfaces.LocalBackup, remoteBackups []interfaces.RemoteBackup) (interfaces.Worker, error) {
	dc, err := helper.GetDockerClient()
	if err != nil {
		return nil, err
	}

	container, err := dc.GetContainer(containerId)
	if err != nil {
		return nil, err
	}

	w := &Worker{
		dc:              dc,
		sourceContainer: container,
	}

	errs := w.mountLocalBackups(localBackups)
	if len(errs) > 0 {
		fmt.Sprintf("Failed to create output volumes: %v", errs)
	}

	errs = w.mountInputVolumes(container.GetVolumes())
	if errs != nil {
		fmt.Sprintf("Failed to create input volumes: %v", err)
	}

	errs = w.mountRemoteBackups(remoteBackups)
	if errs != nil {
		fmt.Sprintf("Failed to create remote backups: %v", err)
	}

	w.passphrase = passphrase

	err = w.createWorkerContainer()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) mountLocalBackup(backup interfaces.LocalBackup) error {
	volume, err := w.dc.CreateVolume(backup.GetVolumeName())
	if err != nil {
		return err
	}

	volume.SetMountPoint(fmt.Sprintf("%s/%s", outputPath, backup.GetVolumeName()))
	backup.SetVolume(volume)
	w.localBackups = append(w.localBackups, backup)
	w.outputVolumes = append(w.outputVolumes, volume)
	return nil
}

func (w *Worker) mountLocalBackups(backups []interfaces.LocalBackup) []error {
	var errs []error
	var errsChan = make(chan error, len(backups))
	var wg sync.WaitGroup

	for _, backup := range backups {
		wg.Add(1)

		go func(backup interfaces.LocalBackup) {
			defer wg.Done()

			err := w.mountLocalBackup(backup)
			if err != nil {
				errsChan <- err
			}
		}(backup)
	}

	go func() {
		wg.Wait()
		close(errsChan)
	}()

	for err := range errsChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (w *Worker) mountKeyfile(keyfile string) interfaces.DockerBind {
	keyfileName := strings.Split(keyfile, "/")[len(strings.Split(keyfile, "/"))-1]
	return docker.NewDockerBind(keyfile, fmt.Sprintf("%s/%s", keyfilesPath, keyfileName), false)
}

func (w *Worker) mountRemoteBackup(backup interfaces.RemoteBackup) error {
	// import ssh key
	keyfile := w.mountKeyfile(backup.GetSshKey())
	w.keyfiles = append(w.keyfiles, keyfile)
	w.remoteBackups = append(w.remoteBackups, backup)

	return nil
}

func (w *Worker) mountRemoteBackups(backups []interfaces.RemoteBackup) []error {
	var errs []error
	var errsChan = make(chan error, len(backups))
	var wg sync.WaitGroup

	for _, backup := range backups {
		wg.Add(1)

		go func(backup interfaces.RemoteBackup) {
			defer wg.Done()

			err := w.mountRemoteBackup(backup)
			if err != nil {
				errsChan <- err
			}
		}(backup)
	}

	go func() {
		wg.Wait()
		close(errsChan)
	}()

	for err := range errsChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (w *Worker) mountInputVolume(volume interfaces.DockerVolume) error {
	path := fmt.Sprintf("%s/%s", inputPath, volume.GetName())
	volume.SetMountPoint(path)
	w.inputVolumes = append(w.inputVolumes, volume)

	return nil
}

func (w *Worker) mountInputVolumes(volumes []interfaces.DockerVolume) []error {
	var errs []error
	var errsChan = make(chan error, len(volumes))
	var wg sync.WaitGroup

	for _, volume := range volumes {
		wg.Add(1)

		go func(volume interfaces.DockerVolume) {
			defer wg.Done()

			err := w.mountInputVolume(volume)
			if err != nil {
				errsChan <- err
			}
		}(volume)
	}

	go func() {
		wg.Wait()
		close(errsChan)
	}()

	for err := range errsChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (w *Worker) createWorkerContainer() error {
	// create worker container
	var err error
	w.workerContainer, err = w.dc.CreateContainer(DOCKER_IMAGE, append(w.inputVolumes, w.outputVolumes...), w.keyfiles)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) initSSHClient() []error {
	var hosts []string
	for _, backupLocation := range w.remoteBackups {
		hosts = append(hosts, backupLocation.GetHost())
	}

	ssh, errs := ssh.NewSSHClient(w.workerContainer, w.keyfiles, hosts)
	if len(errs) > 0 {
		return errs
	}

	w.ssh = ssh
	return nil
}

func (w *Worker) Backup() error {
	if w.workerContainer == nil {
		err := w.createWorkerContainer()
		if err != nil {
			return err
		}
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

	defer w.workerContainer.StopAndRemove()

	bc, err := borg.NewBorgClient(w.workerContainer)
	if err != nil {
		return err
	}

	w.bc = bc

	for _, backupLocation := range w.localBackups {
		path := fmt.Sprintf("%s/%s/%s", outputPath, backupLocation.GetVolumeName(), w.sourceContainer.GetID())
		volumeErrs := w.backupVolumes(w.inputVolumes, path)

		if len(volumeErrs) > 0 {
			fmt.Errorf("Failed to backup volumes: %v", volumeErrs)
		}
	}

	if len(w.remoteBackups) > 0 {
		errs := w.initSSHClient()
		if len(errs) > 0 {
			fmt.Printf("Failed to initialize SSH client: %v", errs)
		}

		var remoteBackupErros [][]error

		for _, backupLocation := range w.remoteBackups {
			path := fmt.Sprintf("%s@%s:%s/%s", backupLocation.GetUser(), backupLocation.GetHost(), backupLocation.GetPath(), w.sourceContainer.GetID())
			volumeErrs := w.backupVolumes(w.inputVolumes, path)
			remoteBackupErros = append(remoteBackupErros, volumeErrs)
		}

		if len(errs) > 0 {
			fmt.Errorf("Failed to backup volumes: %v", errs)
		}
	}

	// post-backup

	return nil
}

func (w *Worker) backupVolumes(volumes []interfaces.DockerVolume, output string) []error {
	var errs []error
	for _, volume := range volumes {
		err := w.backupVolume(volume, fmt.Sprintf("%s/%s", output, volume.GetName()))
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (w *Worker) backupVolume(volume interfaces.DockerVolume, repoPath string) error {
	// create repository
	repo, err := w.bc.GetOrCreateRepo(interfaces.CreateBorgRepoConfig{
		Path:           repoPath,
		Passphrase:     w.passphrase,
		EncryptionType: "none",
		MakeParentDirs: true,
	})

	if err != nil {
		return err
	}

	// backup volume
	err = repo.Backup(volume.GetMountPoint())
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Stop() error {
	if w.workerContainer != nil {
		return w.workerContainer.StopAndRemove()
	}

	return nil
}
