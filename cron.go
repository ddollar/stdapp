package stdapp

import (
	"context"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ddollar/stdcli"
	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	"github.com/kballard/go-shellquote"
	"github.com/robfig/cron/v3"
)

const cronPrefix = "stdapp.cron."

var cronEntryMatchers = []*regexp.Regexp{
	regexp.MustCompile(`^(@every (?:\d+[smhdw])+) (.*)$`),
	regexp.MustCompile(`^(@(?:annually|yearly|monthly|weekly|daily|hourly|reboot)) (.*)$`),
	regexp.MustCompile(`^((?:(?:(?:\d+,)+\d+|(?:\d+(?:\/|-|#)\d+)|\d+L?|\*(?:\/\d+)?|L(?:-\d+)?|\?|[A-Z]{3}(?:-[A-Z]{3})?) ?){5,6}) (.*)$`),
}

type Cron struct {
	ctx    *stdcli.Context
	docker *docker.Client
}

func NewCron(ctx *stdcli.Context) (*Cron, error) {
	dc, err := dockerClient()
	if err != nil {
		return nil, err
	}

	c := &Cron{
		ctx:    ctx,
		docker: dc,
	}

	return c, nil
}

func (cc *Cron) Run() error {
	cs, err := dockerProjectContainers(cc.docker)
	if err != nil {
		return err
	}

	cr := cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))

	for _, c := range cs {
		for k, v := range c.Labels {
			if !strings.HasPrefix(k, cronPrefix) {
				continue
			}

			if schedule, command, ok := cronEntry(v); ok {
				name := strings.TrimPrefix(k, cronPrefix)
				cc.ctx.Writef("ns=cron at=start name=%q schedule=%q command=%q\n", name, schedule, command)

				if _, err := cr.AddFunc(schedule, cc.job(name, c.ID, command)); err != nil {
					return err
				}
			}
		}
	}

	cr.Run()

	return nil
}

func (cc *Cron) exec(id, command string) error {
	ctx := context.Background()

	cmd, err := shellquote.Split(command)
	if err != nil {
		return err
	}

	copts := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	e, err := cc.docker.ContainerExecCreate(ctx, id, copts)
	if err != nil {
		return err
	}

	sopts := types.ExecStartCheck{
		// OutputStream: os.Stdout,
		// ErrorStream:  os.Stderr,
	}

	res, err := cc.docker.ContainerExecAttach(ctx, e.ID, sopts)
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, res.Reader)

	return nil
}

func (cc *Cron) job(name, id, command string) func() {
	return func() {
		cc.ctx.Writef("ns=cron at=run name=%q command=%q\n", name, command)

		if err := cc.exec(id, command); err != nil {
			cc.ctx.Writef("ns=cron at=run name=%q error=%q\n", name, err)
			return
		}
	}
}

func cronEntry(entry string) (string, string, bool) {
	for _, m := range cronEntryMatchers {
		if parts := m.FindStringSubmatch(entry); len(parts) == 3 {
			return parts[1], parts[2], true
		}
	}

	return "", "", false
}
