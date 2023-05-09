package migrate

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/go-pg/pg/v10"
)

func Run(dburl string, migrations fs.FS) error {
	opts, err := pg.ParseURL(dburl)
	if err != nil {
		return err
	}

	db := pg.Connect(opts)

	e := &Engine{
		db: db,
		fs: migrations,
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
