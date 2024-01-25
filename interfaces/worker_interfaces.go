package interfaces

type Worker interface {
	Exec(cmd string) (string, error)
	GetBorgClient() BorgClient
	GetEnv(key string) string
	GetEnvs() map[string]string
	SetEnv(key string, value string)
	StopAndRemove() error
	Backup() error
}
