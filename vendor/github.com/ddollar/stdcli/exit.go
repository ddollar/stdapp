package stdcli

import "fmt"

type ExitCoder interface {
	ExitCode() int
}

type exitCode struct {
	code int
}

func (e exitCode) ExitCode() int {
	return e.code
}

func (e exitCode) Error() string {
	return fmt.Sprintf("exit %d", e.code)
}

func Exit(code int) ExitCoder {
	return exitCode{code: code}
}
