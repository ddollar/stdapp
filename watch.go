package stdapp

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/ddollar/errors"
	"github.com/fsnotify/fsnotify"
)

const debounce = 100 * time.Millisecond

func (a *App) spawn(command string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command("go", append([]string{"run", ".", command}, args...)...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err)
	}

	return cmd, nil
}

func (a *App) watchAndReload(extensions []string, command string, args ...string) error {
	if len(extensions) == 0 {
		a.logger.At("spawn").Logf("extensions=%q", strings.Join(extensions, ","))

		cmd, err := a.spawn(command, args...)
		if err != nil {
			return errors.Wrap(err)
		}

		if err := cmd.Wait(); err != nil {
			return errors.Wrap(err)
		}

		return nil
	}

	a.logger.At("watch").Logf("extensions=%q", strings.Join(extensions, ","))

	ch := make(chan string)

	if err := a.watchChanges(extensions, ch); err != nil {
		return errors.Wrap(err)
	}

	for {
		cmd, err := a.spawn(command, args...)
		if err != nil {
			return errors.Wrap(err)
		}

		a.logger.At("change").Logf("file=%q", <-ch)

		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
			return errors.Wrap(err)
		}

		if _, err := cmd.Process.Wait(); err != nil {
			return errors.Wrap(err)
		}
	}
}

func (a *App) watchChanges(extensions []string, ch chan<- string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err)
	}

	eh := map[string]bool{}

	for _, e := range extensions {
		eh[e] = true
	}

	paths := map[string]bool{}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		switch path {
		case "vendor":
			return filepath.SkipDir
		}

		if eh[strings.TrimPrefix(filepath.Ext(path), ".")] {
			dir, _ := filepath.Split(path)
			paths[filepath.Join(".", dir)] = true
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err)
	}

	for path := range paths {
		if err := w.Add(path); err != nil {
			return errors.Wrap(err)
		}
	}

	t := time.NewTimer(1 * time.Hour)

	if !t.Stop() {
		<-t.C
	}

	go watchLoop(ch, w, t)

	return nil
}

func watchLoop(ch chan<- string, w *fsnotify.Watcher, t *time.Timer) {
	var e fsnotify.Event

	for {
		select {
		case e = <-w.Events:
			if e.Op.Has(fsnotify.Write) || e.Op.Has(fsnotify.Chmod) {
				t.Reset(debounce)
			}
		case <-t.C:
			ch <- e.Name
		}
	}
}
