package main

import (
	"docker-backup/internal/helper"
	"docker-backup/internal/worker"
	"fmt"
)

const (
	targetContainer = "test-service"
	configPath      = "/Users/jannis/Git/docker-backup/docker-backup.yml"
)

func backupContainer() {
	localBackups, remoteBackups, err := helper.ParseConfigFile(configPath)
	if err != nil {
		fmt.Printf("error parsing config file: %s\n", err)
		return
	}

	fmt.Printf("localBackups: %v\n", localBackups)
	fmt.Printf("remoteBackups: %v\n", remoteBackups)

	w, err := worker.NewWorker(targetContainer, localBackups, remoteBackups)
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
