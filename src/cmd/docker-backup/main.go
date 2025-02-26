package main

import (
	"context"
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/db"
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

func getTargetContainers(project utils.Project, cli *docker.Client) []*docker.Container {
	var containers []*docker.Container
	for _, id := range project.Containers {
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
	config, err := utils.ParseConfig()
	ctx := context.Background()

	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	dbClient := db.NewClient()
	fmt.Sprintf("Connected to database %v", dbClient)

	client, err := docker.NewClient(ctx, config.WorkerImage)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	defer client.Close()

	projects := config.Projects
	if len(projects) == 0 {
		log.Fatal("No projects found")
	}

	for name, project := range projects {
		p, err := dbClient.GetProject(name)
		if err != nil {
			if err.Error() == "record not found" {
				fmt.Sprintf("Project %s not found, creating...", name)
				p, err = dbClient.CreateProject(name)
				if err != nil {
					log.Printf("Error creating project %s: %v", name, err)
					continue
				}

				log.Printf("Project %s created.", name)
			} else {
				fmt.Errorf("Error getting project %s: %v", name, err)
				continue
			}
		}

		containers := getTargetContainers(project, client)
		if len(containers) == 0 {
			log.Printf("No containers found for project %s", name)
			continue
		}

		mounts := getInputVolumes(containers)
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: config.Volume,
			Target: "/output",
		})

		worker, err := client.SpawnWorker(name, project.Passphrase, mounts)
		if err != nil {
			fmt.Errorf("Error spawning worker: %v", err)
			continue
		}

		defer func() {
			if err := worker.StopAndRemove(ctx); err != nil {
				log.Printf("Error removing worker container: %v", err)
			} else {
				log.Printf("Worker container %s removed.", worker.ID)
			}
		}()

		if err := worker.InitOutputRepo(); err != nil {
			log.Printf("Error initializing output repository: %v", err)
			worker.StopAndRemove(ctx)
			continue
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
			continue
		}

		for _, c := range pausedContainers {
			if err := c.Unpause(ctx); err != nil {
				log.Printf("Error unpausing container %s: %v", c.ID, err)
			}
		}

		remotes := config.Remotes
		for name, remote := range remotes {
			if err := worker.Sync(name, remote); err != nil {
				log.Printf("Error syncing to remote %s: %v", name, err)
				continue
			}
		}

		_, err = p.CreateBackup()
		if err != nil {
			log.Printf("Error storing backup in db %s: %v", name, err)
		}

		log.Printf("Backup for project %s completed.", name)
	}

	log.Println("Backup process completed.")
}
