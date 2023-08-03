package migrate

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

type Migration struct {
	Version string
	Body    []byte
}

type Migrations []Migration

func LoadMigrations(e *Engine) (Migrations, error) {
	raw := map[string]Migration{}

	root := "."

	if e.dir != "" {
		root = e.dir
	}

	err := fs.WalkDir(e.fs, root, func(path string, d fs.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		parts := strings.SplitN(d.Name(), ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid migration: %s", d.Name())
		}
		mm := raw[parts[0]]
		mm.Body, err = fs.ReadFile(e.fs, path)
		if err != nil {
			return err
		}
		raw[parts[0]] = mm
		return nil
	})
	if err != nil {
		return nil, err
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
