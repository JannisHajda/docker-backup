package db

type Project struct {
	db         *Database
	ID         int
	Containers []*Container
	Name       string
}

func (p *Project) AddContainer(containerID string) error {
	err := p.db.pct.add(p.ID, containerID)
	if err != nil {
		return err
	}

	return nil
}
