package stdapp

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

const debounce = 100 * time.Millisecond

func (a *App) watchAndReload(extensions []string, cmd string, args ...string) error {
	a.logger.At("watch").Logf("extensions=%q", strings.Join(extensions, ","))

	ch := make(chan string)

	if err := a.watchChanges(extensions, ch); err != nil {
		return err
	}

	for {
		cmd := exec.Command("go", append([]string{"run", ".", cmd}, args...)...)

		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return errors.WithStack(err)
		}

		a.logger.At("change").Logf("file=%q", <-ch)

		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
			return errors.WithStack(err)
		}

		if _, err := cmd.Process.Wait(); err != nil {
			return errors.WithStack(err)
		}
	}
}

func (a *App) watchChanges(extensions []string, ch chan<- string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}

	for path := range paths {
		if err := w.Add(path); err != nil {
			return errors.WithStack(err)
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
