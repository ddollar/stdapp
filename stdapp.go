package stdapp

import (
	"fmt"

	"github.com/ddollar/logger"
	"github.com/pkg/errors"
)

func New(opts Options) (*App, error) {
	// if opts.Name == "" {
	// 	return nil, errors.WithStack(errors.New("Name required"))
	// }

	// if opts.Resolver == nil {
	// 	return nil, errors.WithStack(errors.New("Resolver required"))
	// }

	// if opts.Schema == "" {
	// 	return nil, errors.WithStack(errors.New("Schema required"))
	// }

	a := &App{
		compose:    opts.Compose,
		database:   opts.Database,
		logger:     logger.New(fmt.Sprintf("ns=%s", opts.Name)),
		migrations: opts.Migrations,
		name:       opts.Name,
		schema:     opts.Schema,
	}

	if opts.Database != "" && opts.Resolver != nil {
		db, err := database(opts.Database)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		r, err := opts.Resolver(db)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		a.resolver = r
	}

	return a, nil
}
