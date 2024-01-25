package interfaces

type Worker interface {
	BackupContainer(containerIdentifier string) error
	BackupProject(projectName string) error
}
