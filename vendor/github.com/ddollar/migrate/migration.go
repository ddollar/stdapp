package migrate

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/ddollar/errors"
)

type Migration struct {
	Version string
	Body    []byte
}

type Migrations []Migration

func LoadMigrations(e *Engine) (Migrations, error) {
	raw := map[string]Migration{}

	files, err := fs.ReadDir(e.fs, ".")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		parts := strings.SplitN(file.Name(), ".", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid migration: %s", file.Name())
		}

		mm := raw[parts[0]]

		mm.Body, err = fs.ReadFile(e.fs, file.Name())
		if err != nil {
			return nil, err
		}

		raw[parts[0]] = mm
	}

	ms := Migrations{}

	for k, m := range raw {
		m.Version = k
		ms = append(ms, m)
	}

	sort.Slice(ms, func(i, j int) bool { return ms[i].Version < ms[j].Version })

	return ms, nil
}

func (ms Migrations) Find(version string) (Migration, bool) {
	for _, m := range ms {
		if m.Version == version {
			return m, true
		}
	}

	return Migration{}, false
}
