package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"net/url"
)

type Options struct {
	Dir    string
	DryRun bool
}

func Run(ctx context.Context, dburl string, migrations fs.FS, opts Options) error {
	u, err := url.Parse(dburl)
	if err != nil {
		return err
	}

	db, err := sql.Open(u.Scheme, dburl)
	if err != nil {
		return err
	}

	e := &Engine{
		db:     db,
		dir:    opts.Dir,
		dryrun: opts.DryRun,
		fs:     migrations,
	}

	if err := e.Initialize(); err != nil {
		return err
	}

	ms, err := e.Pending()
	if err != nil {
		return err
	}

	for _, m := range ms {
		fmt.Printf("%s: ", m)

		if err := e.Migrate(ctx, m); err != nil {
			fmt.Printf("%s\n", err)
			return errors.New("migration failed")
		} else {
			fmt.Println("OK")
		}
	}

	return nil
}
