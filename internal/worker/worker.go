package worker

import (
	"docker-backup/errors"
	"docker-backup/interfaces"
	"docker-backup/internal/borg"
	"docker-backup/internal/docker"
	"docker-backup/internal/helper"
	"docker-backup/internal/ssh"
	goerrors "errors"
	"fmt"
	"strings"
	"time"
)

const (
	WORKER_IMAGE      = "worker"
	inputVolumesPath  = "/input"
	outputVolumesPath = "/output"
	sshKeyfilesPath   = "/ssh_keyfiles"
	borgKeyfilesPath  = "/borg_keyfiles"
)

type Worker struct {
	db  interfaces.DatabaseClient
	dc  interfaces.DockerClient
	bc  interfaces.BorgClient
	ssh interfaces.SSHClient

	sourceContainer interfaces.DockerContainer
	workerContainer interfaces.DockerContainer

	inputVolumes []interfaces.DockerVolume
	borgKeyfiles []interfaces.DockerBind
	sshKeyfiles  []interfaces.DockerBind

	localBackups  []interfaces.LocalBackup
	remoteBackups []interfaces.RemoteBackup
}

func NewWorker(containerId string, localBackups []interfaces.LocalBackup, remoteBackups []interfaces.RemoteBackup) (interfaces.Worker, error) {
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

	var inputVolumes []interfaces.DockerVolume
	var outputVolumes []interfaces.DockerVolume
	var borgKeyfiles []interfaces.DockerBind
	var sshKeyfiles []interfaces.DockerBind
	var hosts []string

	for _, v := range container.GetVolumes() {
		path := fmt.Sprintf("%s/%s", inputVolumesPath, v.GetName())
		v.SetMountPoint(path)
		inputVolumes = append(inputVolumes, v)
	}

	for _, backup := range localBackups {
		v, err := dc.CreateVolume(backup.VolumeName)
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf("%s/%s", outputVolumesPath, v.GetName())
		v.SetMountPoint(path)

		borgKeyfile := backup.Keyfile

		// mount the borg keyfile
		if borgKeyfile != "" {
			bind := docker.NewDockerBind(borgKeyfile, fmt.Sprintf("%s/%s", borgKeyfilesPath, borgKeyfile), false)
			borgKeyfiles = append(borgKeyfiles, bind)
		}

		outputVolumes = append(outputVolumes, v)
	}

	for _, backup := range remoteBackups {
		borgKeyfile := backup.Keyfile
		sshKeyfile := backup.SSHKey

		// mount the borg keyfile
		if borgKeyfile != "" {
			bind := docker.NewDockerBind(borgKeyfile, fmt.Sprintf("%s/%s", borgKeyfilesPath, borgKeyfile), false)
			borgKeyfiles = append(borgKeyfiles, bind)
		}

		// mount the ssh keyfile
		if sshKeyfile != "" {
			keyfileName := strings.Split(sshKeyfile, "/")[len(strings.Split(sshKeyfile, "/"))-1]
			bind := docker.NewDockerBind(sshKeyfile, fmt.Sprintf("%s/%s", sshKeyfilesPath, keyfileName), false)
			sshKeyfiles = append(sshKeyfiles, bind)
		}

		hosts = append(hosts, backup.Host)
	}

	workerContainer, err := dc.CreateContainer(WORKER_IMAGE, append(inputVolumes, outputVolumes...), append(borgKeyfiles, sshKeyfiles...))
	if err != nil {
		return nil, err
	}

	workerContainer.Start()

	sshClient, errs := ssh.NewSSHClient(workerContainer, sshKeyfiles, hosts)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	w.ssh = sshClient

	borgClient, err := borg.NewBorgClient(workerContainer)
	if err != nil {
		return nil, err
	}

	w.bc = borgClient
	w.inputVolumes = inputVolumes
	w.workerContainer = workerContainer
	w.localBackups = localBackups
	w.remoteBackups = remoteBackups

	return w, nil
}

func (w *Worker) Stop() error {
	err := w.workerContainer.StopAndRemove()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Backup() error {
	var sources []string
	for _, v := range w.inputVolumes {
		sources = append(sources, v.GetMountPoint())
	}

	var localBackupErrors []error
	for _, localBackup := range w.localBackups {
		repoPath := fmt.Sprintf("%s/%s/%s", outputVolumesPath, localBackup.VolumeName, w.sourceContainer.GetID())
		repo, err := w.createOrGetRepo(repoPath, localBackup.Passphrase)
		if err != nil {
			localBackupErrors = append(localBackupErrors, err)
			continue
		}

		currentTimestamp := time.Now().Format("2006-01-02T15:04:05")
		err = repo.CreateArchive(interfaces.CreateBorgArchiveConfig{
			Compression: "zstd",
			Sources:     sources,
			Name:        currentTimestamp,
		})

		if err != nil {
			localBackupErrors = append(localBackupErrors, err)
		}
	}

	var remoteBackupErrors []error
	for _, remoteBackup := range w.remoteBackups {
		repoPath := fmt.Sprintf("%s@%s:%s/%s", remoteBackup.User, remoteBackup.Host, remoteBackup.Path, w.sourceContainer.GetID())
		repo, err := w.createOrGetRepo(repoPath, remoteBackup.Passphrase)
		if err != nil {
			remoteBackupErrors = append(remoteBackupErrors, err)
			continue
		}

		currentTimestamp := time.Now().Format("2006-01-02T15:04:05")

		err = repo.CreateArchive(interfaces.CreateBorgArchiveConfig{
			Compression: "zstd",
			Sources:     sources,
			Name:        currentTimestamp,
		})

		if err != nil {
			remoteBackupErrors = append(remoteBackupErrors, err)
		}
	}

	if len(localBackupErrors) > 0 {
		fmt.Printf("Local backup errors: %v\n", localBackupErrors)
	}

	if len(remoteBackupErrors) > 0 {
		fmt.Printf("Remote backup errors: %v\n", remoteBackupErrors)
	}

	return nil
}

func (w *Worker) createOrGetRepo(repoPath string, passphrase string) (interfaces.BorgRepo, error) {
	repo, err := w.bc.GetRepo(interfaces.GetBorgRepoConfig{Path: repoPath, Passphrase: passphrase})
	if err != nil {
		var err *errors.RepositoryDoesNotExistError
		if goerrors.As(err, &err) {
			var createErr error
			repo, createErr = w.bc.CreateRepo(interfaces.CreateBorgRepoConfig{Path: repoPath, EncryptionType: "repokey-blake2", Passphrase: passphrase})

			if createErr != nil {
				return nil, err
			}
		}
	}

	return repo, nil
}
