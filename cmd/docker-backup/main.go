package main

import (
	"context"
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/docker"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	_ "github.com/JannisHajda/docker-backup/internal/docker"
	"github.com/docker/docker/api/types/mount"
)

type Env struct {
	TARGET_CONTAINERS []string
	BORG_REPO         string
	BORG_PASSPHRASE   string
	OUTPUT_VOLUME     string
	MEGA_USERNAME     string
	MEGA_PASSWORD     string
}

func loadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	env := Env{
		BORG_REPO:       os.Getenv("BORG_REPO"),
		BORG_PASSPHRASE: os.Getenv("BORG_PASSPHRASE"),
		OUTPUT_VOLUME:   os.Getenv("OUTPUT_VOLUME"),
		MEGA_USERNAME:   os.Getenv("MEGA_USERNAME"),
		MEGA_PASSWORD:   os.Getenv("MEGA_PASSWORD"),
	}

	targetContainers := os.Getenv("TARGET_CONTAINERS")
	if targetContainers != "" {
		env.TARGET_CONTAINERS = strings.Split(targetContainers, ",")
	}

	if len(env.TARGET_CONTAINERS) == 0 || env.BORG_REPO == "" || env.BORG_PASSPHRASE == "" || env.OUTPUT_VOLUME == "" || env.MEGA_USERNAME == "" || env.MEGA_PASSWORD == "" {
		log.Fatal("BORG_REPO, BORG_PASSPHRASE, and OUTPUT_VOLUME must be set in .env file")
	}

	return env
}

func getUniqueVolumeMounts(containers []*docker.Container) []mount.Mount {
	volumesMap := make(map[string]bool)
	var mounts []mount.Mount
	for _, c := range containers {
		for _, m := range c.Mounts {
			if m.Type == mount.TypeVolume && !volumesMap[m.Name] {
				volumesMap[m.Name] = true
				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeVolume,
					Source: m.Name,
					Target: fmt.Sprintf("/input/%s", m.Name),
				})
			}
		}
	}

	return mounts
}

func getTargetContainers(env Env, cli *docker.Client) []*docker.Container {
	var containers []*docker.Container
	for _, id := range env.TARGET_CONTAINERS {
		c, err := cli.GetContainer(id)
		if err != nil {
			log.Printf("Error getting container %s: %v", id, err)
			continue
		}

		containers = append(containers, c)
	}

	return containers
}

func main() {
	env := loadEnv()
	ctx := context.Background()

	client, err := docker.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	defer client.Close()

	containers := getTargetContainers(env, client)
	if len(containers) == 0 {
		log.Fatal("No target containers found")
	}

	mounts := getUniqueVolumeMounts(containers)
	mounts = append(mounts, mount.Mount{
		Type:   mount.TypeVolume,
		Source: "backups",
		Target: "/output",
	})

	worker, err := client.SpawnWorker(env.BORG_REPO, env.BORG_PASSPHRASE, mounts)
	if err != nil {
		log.Fatalf("Error spawning worker container: %v", err)
	}

	defer func() {
		if err := worker.StopAndRemove(ctx); err != nil {
			log.Printf("Error removing worker container: %v", err)
		} else {
			log.Printf("Worker container %s removed.", worker.ID)
		}
	}()

	if err := worker.InitOutputRepo(); err != nil {
		log.Fatalf("Error initializing output repository: %v", err)
		worker.StopAndRemove(ctx)
		return
	}

	var pausedContainers []*docker.Container
	for _, c := range containers {
		labels := c.Config.Labels
		preBackupCommand := labels["docker-backup.prebackup"]

		if preBackupCommand != "" {
			_, stderr, exitCode, err := c.Exec(preBackupCommand)
			if err != nil {
				log.Printf("Error running pre-backup command for container %s: %v", c.ID, err)
				continue
			} else if exitCode != 0 {
				log.Printf("Pre-backup command failed for container %s: %s", c.ID, stderr)
				continue
			}
		}

		if err := c.Pause(ctx); err != nil {
			log.Printf("Error pausing container %s: %v", c.Name, err)
			continue
		}

		pausedContainers = append(pausedContainers, c)
	}

	if err := worker.BackupRepo(); err != nil {
		log.Printf("Error backing up repository: %v", err)
		worker.StopAndRemove(ctx)
		return
	}

	for _, c := range pausedContainers {
		if err := c.Unpause(ctx); err != nil {
			log.Printf("Error unpausing container %s: %v", c.ID, err)
		}
	}

	log.Println("Backup process completed.")
}
