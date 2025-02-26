package interfaces

type Shell interface {
	Exec(cmd string) (stdout string, stderr string, exitCode int, err error)
}
