package stdapp

import "os/exec"

type Runner func(cmd string, args ...string) *exec.Cmd

func RunnerCompose(container string, tty bool) Runner {
	return func(cmd string, args ...string) *exec.Cmd {
		eargs := []string{"compose", "exec"}
		if !tty {
			eargs = append(eargs, "-T")
		}
		eargs = append(eargs, container, cmd)
		return exec.Command("docker", append(eargs, args...)...)
	}
}

func RunnerLocal(cmd string, args ...string) *exec.Cmd {
	return exec.Command(cmd, args...)
}
