package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
)

type Engine struct {
	db         *sql.DB
	dir        string
	dryrun     bool
	fs         fs.FS
	migrations Migrations
	state      State
}

func (e *Engine) Initialize() error {
	if _, err := e.db.Exec(`create table if not exists "_migrations" (version varchar unique not null);`); err != nil {
		return err
	}

	ms, err := LoadMigrations(e)
	if err != nil {
		return err
	}

	e.migrations = ms

	ss, err := LoadState(e)
	if err != nil {
		return err
	}

	e.state = ss

	return nil
}

func (e *Engine) Migrate(ctx context.Context, version string) error {
	m, ok := e.migrations.Find(version)
	if !ok {
		return fmt.Errorf("no such migration: %s", version)
	}

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := tx.Exec("insert into _migrations values (?)", version); err != nil {
		return err
	}

	if _, err := tx.Exec(string(m.Body)); err != nil {
		return err
	}

	if e.dryrun {
		return tx.Rollback()
	}

	return tx.Commit()
}

func (e *Engine) Pending() ([]string, error) {
	ps := []string{}

	for _, m := range e.migrations {
		if !e.state[m.Version] {
			ps = append(ps, m.Version)
		}
	}

	sort.Strings(ps)

	return ps, nil
}
