package main

import (
	"docker-backup/internal/worker"
	"fmt"
)

const (
	targetContainer = "test-service"
	passphrase      = "test2"
	outputVolume    = "test-output"
)

func backupContainer() {
	worker, err := worker.NewWorker(targetContainer, outputVolume, passphrase)
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
