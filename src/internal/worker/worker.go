package worker

import (
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/borg"
	"github.com/JannisHajda/docker-backup/internal/docker"
	"github.com/JannisHajda/docker-backup/internal/rclone"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
)

type Worker struct {
	container *docker.Container
	borg      *borg.Client
	repo      *borg.Repo
	rclone    *rclone.Client
	remotes   []*rclone.Remote
}

type RepoConfig struct {
	Name       string
	Passphrase string
}

func initContainer(dockerClient *docker.Client, image string, mounts []mount.Mount, env []string) (*docker.Container, error) {
	config := container.Config{
		Image: image,
		Cmd:   []string{"tail", "-f", "/dev/null"},
		Env:   env,
	}

	hostConfig := container.HostConfig{
		Mounts: mounts,
	}

	ctx := dockerClient.GetContext()
	resp, err := dockerClient.ContainerCreate(ctx, &config, &hostConfig, nil, nil, "")
	if err != nil {
		return nil, err
	}

	c, err := dockerClient.GetContainer(resp.ID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func NewWorker(dockerClient *docker.Client, image string, mounts []mount.Mount, env []string) (*Worker, error) {
	container, err := initContainer(dockerClient, image, mounts, env)
	if err != nil {
		return nil, err
	}

	if err := container.Start(); err != nil {
		return nil, err
	}

	borgClient, err := borg.NewClient(container)
	if err != nil {
		return nil, err
	}

	rcloneClient, err := rclone.NewClient(container)
	if err != nil {
		return nil, err
	}

	w := Worker{
		container: container,
		borg:      borgClient,
		rclone:    rcloneClient,
	}

	return &w, nil
}

func (w *Worker) Backup(repoConfig RepoConfig) (*borg.Archive, error) {
	repoPath := fmt.Sprintf("/output/%s", repoConfig.Name)
	repo, err := w.borg.GetRepo(repoPath, repoConfig.Passphrase)
	if err != nil {
		if err.Error() != "failed to list borg repository: exit code 2" {
			return nil, err
		}

		repo, err = w.borg.CreateRepo(repoPath, repoConfig.Passphrase)
		if err != nil {
			return nil, err
		}
	}

	archive, err := repo.Backup("/input", repoConfig.Passphrase)
	if err != nil {
		return nil, err
	}

	for _, remote := range w.remotes {
		err := w.rclone.Sync(rclone.SyncConfig{
			InputPath:  repo.Path,
			OutputPath: "/testbackup",
			Remote:     *remote,
		})

		if err != nil {
			return nil, err
		}
	}

	return archive, nil
}

func (w *Worker) Stop() error {
	return w.container.StopAndRemove()
}
