package stdapp

import (
	"fmt"
	"os/exec"
)

type Runner func(cmd string, args ...string) *exec.Cmd

func RunnerCompose(container string, tty bool, env map[string]string) Runner {
	return func(cmd string, args ...string) *exec.Cmd {
		eargs := []string{"compose", "exec"}
		if !tty {
			eargs = append(eargs, "-T")
		}
		for k, v := range env {
			eargs = append(eargs, "--env", fmt.Sprintf("%s=%s", k, v))
		}
		eargs = append(eargs, container, cmd)
		return exec.Command("docker", append(eargs, args...)...)
	}
}

func RunnerLocal(cmd string, args ...string) *exec.Cmd {
	return exec.Command(cmd, args...)
}
