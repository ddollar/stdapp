package migrate

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/go-pg/pg/v10"
)

type Options struct {
	Dir    string
	Schema string
}

func Run(dburl string, migrations fs.FS, opts Options) error {
	dbopts, err := pg.ParseURL(dburl)
	if err != nil {
		return err
	}

	if opts.Schema != "" {
		dbopts.OnConnect = func(ctx context.Context, conn *pg.Conn) error {
			_, err := conn.Exec("SET search_path=?", opts.Schema)
			return err
		}
	}

	db := pg.Connect(dbopts)

	e := &Engine{
		db:  db,
		dir: opts.Dir,
		fs:  migrations,
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

		if err := e.Migrate(m); err != nil {
			fmt.Printf("%s\n", err)
			return errors.New("migration failed")
		} else {
			fmt.Println("OK")
		}
	}

	return nil
}
