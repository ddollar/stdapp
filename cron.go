package stdapp

import (
	"os"
	"regexp"
	"strings"

	"github.com/ddollar/stdcli"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/kballard/go-shellquote"
	"github.com/robfig/cron/v3"
)

const cronPrefix = "stdapp.cron."

var cronEntryMatchers = []*regexp.Regexp{
	regexp.MustCompile(`^(@every (?:\d+[smhdw])+) (.*)$`),
	regexp.MustCompile(`^(@(?:annually|yearly|monthly|weekly|daily|hourly|reboot)) (.*)$`),
	regexp.MustCompile(`^((?:(?:(?:\d+,)+\d+|(?:\d+(?:\/|-|#)\d+)|\d+L?|\*(?:\/\d+)?|L(?:-\d+)?|\?|[A-Z]{3}(?:-[A-Z]{3})?) ?){5,7}) (.*)$`),
}

func cronStart(ctx *stdcli.Context, dc *docker.Client, cs []docker.APIContainers) error {
	cr := cron.New()

	for _, c := range cs {
		for k, v := range c.Labels {
			if strings.HasPrefix(k, cronPrefix) {
				name := strings.TrimPrefix(k, cronPrefix)

				if schedule, command, ok := cronEntry(v); ok {
					ctx.Writef("ns=cron at=start name=%q schedule=%q command=%q\n", name, schedule, command)

					_, err := cr.AddFunc(schedule, func() {
						ctx.Writef("ns=cron at=run name=%q command=%q\n", name, command)

						cmd, err := shellquote.Split(command)
						if err != nil {
							ctx.Writef("ns=cron at=run name=%q error=%q\n", name, err)
							return
						}

						copts := docker.CreateExecOptions{
							Cmd:          cmd,
							Container:    c.ID,
							AttachStdout: true,
							AttachStderr: true,
						}

						e, err := dc.CreateExec(copts)
						if err != nil {
							ctx.Writef("ns=cron at=run name=%q error=%q\n", name, err)
							return
						}

						sopts := docker.StartExecOptions{
							OutputStream: os.Stdout,
							ErrorStream:  os.Stderr,
						}

						if err := dc.StartExec(e.ID, sopts); err != nil {
							ctx.Writef("ns=cron at=run name=%q error=%q\n", name, err)
							return
						}
					})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	cr.Run()

	return nil
}

func cronEntry(entry string) (string, string, bool) {
	for _, m := range cronEntryMatchers {
		if parts := m.FindStringSubmatch(entry); len(parts) == 3 {
			return parts[1], parts[2], true
		}
	}

	return "", "", false
}
