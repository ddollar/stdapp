package stdapp

import "os/exec"

type Compose struct {
	Name     string
	Services map[string]struct {
		Environment map[string]string
	}
}

func (c *Compose) Runner(container string) Runner {
	return func(cmd string, args ...string) *exec.Cmd {
		return exec.Command("docker", append([]string{"compose", "exec", "-it", container, cmd}, args...)...)
	}
}
