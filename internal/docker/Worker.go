package docker

import (
	"fmt"
	"github.com/JannisHajda/docker-backup/internal/utils"
	_ "github.com/rclone/rclone/backend/mega" // Import the mega backend
)

type Worker struct {
	Container
	repoName string
}

func (w *Worker) initRepo() error {
	_, _, exitCode, err := w.Exec("borg init --encryption=repokey")
	if err != nil {
		return err
	}

	if exitCode != 0 {
		if exitCode == 2 {

		}

		return fmt.Errorf("failed to initialize borg repository: exit code %d", exitCode)
	}

	return nil
}

func (w *Worker) getRepo() error {
	_, _, exitCode, err := w.Exec("borg list")
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("failed to list borg repository: exit code %d", exitCode)
	}

	return nil
}

func (w *Worker) InitOutputRepo() error {
	if err := w.getRepo(); err != nil {
		if err.Error() != "failed to list borg repository: exit code 2" {
			return err
		}

		if err := w.initRepo(); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) BackupRepo() error {
	cmd := fmt.Sprintf("borg create --stats ::%s /input", "-{now:%Y-%m-%d_%H:%M:%S}")
	_, stderr, exitCode, err := w.Exec(cmd)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("failed to create borg archive: exit code %d, stderr: %s", exitCode, stderr)
	}

	return nil
}

func (w *Worker) Sync(name string, conf utils.Remote) error {
	cmd := fmt.Sprintf("rclone config create %s %s user %s pass %s", name, conf.Type, conf.User, conf.Pass)
	_, stderr, exitCode, err := w.Exec(cmd)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("failed to create rclone config: exit code %d, stderr: %s", exitCode, stderr)
	}

	cmd = fmt.Sprintf("rclone sync /output/%s %s:%s/%s", w.repoName, name, conf.Path, w.repoName)
	_, stderr, exitCode, err = w.Exec(cmd)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("failed to sync rclone: exit code %d, stderr: %s", exitCode, stderr)
	}

	return nil
}
