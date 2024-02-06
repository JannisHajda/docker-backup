package main

import (
	"docker-backup/internal/worker"
	"fmt"
)

func main() {
	worker, err := worker.NewWorker("docker-backup-output")
	if err != nil {
		fmt.Printf("Could not create worker: %s\n", err.Error())
		return
	}

	err = worker.BackupContainer("test-service")
	if err != nil {
		fmt.Printf("Could not backup container: %s\n", err.Error())
		return
	}
}
