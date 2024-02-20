package main

import (
	"docker-backup/interfaces"
	"docker-backup/internal/worker"
	"fmt"
)

const (
	targetContainer = "test-service"
)

func backupContainer() {
	localbackup1 := interfaces.LocalBackup{Backup: interfaces.Backup{
		Passphrase: "test",
	}, VolumeName: "local-backup"}

	remotebackup1 := interfaces.RemoteBackup{Backup: interfaces.Backup{
		Passphrase: "test",
		Path:       "/home/borg/backup",
	}, Host: "remote-backup", User: "borg", SSHKey: "/Users/jannis/Git/docker-backup/build/docker/.ssh/id_rsa"}

	w, err := worker.NewWorker(targetContainer, []interfaces.LocalBackup{localbackup1}, []interfaces.RemoteBackup{remotebackup1})
	if err != nil {
		fmt.Printf("error creating worker: %s\n", err)
		return
	}

	defer w.Stop()

	err = w.Backup()
	if err != nil {
		fmt.Printf("error backing up: %s\n", err)
		return
	}
}

func main() {
	backupContainer()
	//cli.Execute()
}
