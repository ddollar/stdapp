package stdapp

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

const debounce = 100 * time.Millisecond

// test
func (a *App) watchAndExit(extensions []string) {
	if err := a.watchChanges(extensions); err != nil {
		a.logger.At("watch").Logf("error=%q", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func (a *App) watchChanges(extensions []string) error {
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

	var e fsnotify.Event

	t := time.NewTimer(1 * time.Hour)

	if !t.Stop() {
		<-t.C
	}

	for {
		select {
		case e = <-w.Events:
			if e.Op.Has(fsnotify.Write) {
				t.Reset(debounce)
			}
		case <-t.C:
			a.logger.At("change").Logf("file=%q", e.Name)
			return nil
		}
	}
}
