package main

import (
	"docker-backup/internal/worker"
	"fmt"
)

func main() {
	worker, err := worker.NewWorker("test-service", "backups")
	if err != nil {
		fmt.Printf("Could not create worker: %s\n", err.Error())
		return
	}

	err = worker.Backup()
	if err != nil {
		fmt.Printf("Could not backup: %s\n", err.Error())
	}
}
