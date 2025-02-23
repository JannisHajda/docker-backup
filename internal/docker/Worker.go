package docker

import (
	"fmt"
	_ "github.com/rclone/rclone/backend/mega" // Import the mega backend
)

type Worker struct {
	Container
	repoPath       string
	repoPassphrase string
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
