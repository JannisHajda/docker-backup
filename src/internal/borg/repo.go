package borg

import "fmt"

type Repo struct {
	client   *Client
	Name     string
	Path     string
	Archives []*Archive
}

func (r *Repo) Backup(inputPath string, pass string) (*Archive, error) {
	archiveName := fmt.Sprintf("backup-%d", len(r.Archives)+1)
	cmd := fmt.Sprintf("borg create --stats %s::%s %s", r.Path, archiveName, inputPath)
	cmd = r.client.appendAuth(cmd, pass)
	_, stdError, exitCode, err := r.client.shell.Exec(cmd)
	if err != nil {
		return nil, err
	}

	if exitCode != 0 {
		return nil, fmt.Errorf("failed to create borg archive: exit code %d, stderr: %s", exitCode, stdError)
	}

	archive := &Archive{repo: r, name: archiveName}
	r.Archives = append(r.Archives, archive)

	return archive, nil
}
