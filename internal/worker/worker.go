package worker

import (
	"docker-backup/interfaces"
	"docker-backup/internal/borgclient"
	"fmt"
	"regexp"
	"strings"
)

type Worker struct {
	db      interfaces.DatabaseClient
	dc      interfaces.DockerClient
	c       interfaces.DockerContainer
	bc      interfaces.BorgClient
	project interfaces.DatabaseProject
}

func NewWorker(db interfaces.DatabaseClient, dc interfaces.DockerClient, project interfaces.DatabaseProject) (interfaces.Worker, error) {
	w := &Worker{db: db, dc: dc, project: project}

	containers, err := project.GetContainers()
	if err != nil {
		return nil, err
	}

	var dockerContainers []interfaces.DockerContainer
	var errs []error
	for _, container := range containers {
		dockerContainer, err := dc.GetContainer(container.GetID())
		if err != nil {
			errs = append(errs, err)
			continue
		}

		dockerContainers = append(dockerContainers, dockerContainer)
	}

	volumesMap := make(map[string]interfaces.DockerVolume)
	for _, dockerContainer := range dockerContainers {
		for _, dockerVolume := range dockerContainer.GetVolumes() {
			volumesMap[dockerVolume.GetName()] = dockerVolume
		}
	}

	var volumes []interfaces.DockerVolume
	for _, dockerVolume := range volumesMap {
		dockerVolume.SetMountPoint("/input/" + dockerVolume.GetName())
		volumes = append(volumes, dockerVolume)
	}

	w.c, err = dc.CreateContainer("docker-backup", volumes)
	if err != nil {
		return nil, err
	}

	w.bc, err = borgclient.NewBorgClient(w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func extractKeyFromConfig(configContent string) (string, error) {
	// Define a regular expression pattern to match the key value
	re := regexp.MustCompile(`key\s*=\s*([^\n]+)`)

	// Find the matches in the config content
	matches := re.FindStringSubmatch(configContent)

	// Check if a match is found
	if len(matches) < 2 {
		return "", fmt.Errorf("key not found in config content")
	}

	// Extract the key value from the match
	key := strings.TrimSpace(matches[1])

	return key, nil
}

func (w *Worker) Backup() error {
	var errors []error
	volumes := w.c.GetVolumes()
	for _, v := range volumes {
		repo, err := w.bc.GetOrCreateRepository("/output/"+v.GetName(), "test123")
		if err != nil {
			errors = append(errors, err)
			continue
		}

		err = repo.Archive(v.GetMountPoint())
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return errors[0]
	}

	// for each repo: extract key from config file
	for _, v := range volumes {
		output, err := w.c.Exec("cat /output/" + v.GetName() + "/config")
		if err != nil {
			return err
		}

		key, err := extractKeyFromConfig(output)
		if err != nil {
			return err
		}

		fmt.Printf("Key: %s\n", key)
	}

	return nil
}

func (w *Worker) Exec(command string) (string, error) {
	return w.c.Exec(command)
}

func (w *Worker) SetEnv(key, value string) {
	w.c.SetEnv(key, value)
}

func (w *Worker) GetEnv(key string) string {
	return w.c.GetEnv(key)
}

func (w *Worker) GetEnvs() map[string]string {
	return w.c.GetEnvs()
}

func (w *Worker) GetBorgClient() interfaces.BorgClient {
	return w.bc
}

func (w *Worker) StopAndRemove() error {
	return w.c.StopAndRemove()
}
