package main

import (
	"docker-backup/interfaces"
	"docker-backup/internal/worker"
	"fmt"
)

const (
	targetContainer = "test-service"
	passphrase      = "test"
)

func backupContainer() {
	localbackup1 := worker.NewLocalBackup("local-backup1")
	localbackup2 := worker.NewLocalBackup("local-backup2")
	remotebackup := worker.NewRemoteBackup("borg", "remote-backup", "/home/borg/backups", "/Users/jannis/Git/docker-backup/.ssh/id_ed25519")

	worker, err := worker.NewWorker(targetContainer, passphrase, []interfaces.LocalBackup{localbackup1, localbackup2}, []interfaces.RemoteBackup{remotebackup})
	if err != nil {
		fmt.Printf("error creating worker: %s\n", err)
		return
	}

	defer worker.Stop()

	err = worker.Backup()
	if err != nil {
		fmt.Printf("error backing up: %s\n", err)
		return
	}
}

func main() {
	backupContainer()
	//cli.Execute()
}
