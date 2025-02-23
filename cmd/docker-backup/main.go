package main

import (
	"context"
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/docker"
	"github.com/JannisHajda/docker-backup/internal/utils"
	"log"

	_ "github.com/JannisHajda/docker-backup/internal/docker"
	"github.com/docker/docker/api/types/mount"
)

func getInputVolumes(containers []*docker.Container) []mount.Mount {
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

	config, err := utils.ParseConfig()
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	fmt.Println(config)

	client, err := docker.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	defer client.Close()

	containers := getTargetContainers(env, client)
	if len(containers) == 0 {
		log.Fatal("No target containers found")
	}

	mounts := getInputVolumes(containers)
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

	syncConf := docker.SyncConfig{
		Name:       "mega",
		Type:       "mega",
		OutputPath: "/test",
		User:       env.MEGA_USERNAME,
		Password:   env.MEGA_PASSWORD,
	}

	if err := worker.Sync(syncConf); err != nil {
		log.Printf("Error syncing repository: %v", err)
		worker.StopAndRemove(ctx)
		return
	}

	log.Println("Backup process completed.")
}
