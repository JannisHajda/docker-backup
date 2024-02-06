package interfaces

type Worker interface {
	Backup() error
	Stop() error
}
