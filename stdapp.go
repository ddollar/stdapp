package stdapp

import (
	"fmt"

	"github.com/ddollar/logger"
	"github.com/pkg/errors"
)

func New(opts Options) (*App, error) {
	if opts.Name == "" {
		return nil, errors.New("Name required")
	}

	if opts.Resolver == nil {
		return nil, errors.New("Resolver required")
	}

	if opts.Schema == "" {
		return nil, errors.New("Schema required")
	}

	db, err := database(opts.Database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r, err := opts.Resolver(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	a := &App{
		database:   db,
		logger:     logger.New(fmt.Sprintf("ns=%s", opts.Name)),
		migrations: opts.Migrations,
		name:       opts.Name,
		resolver:   r,
		schema:     opts.Schema,
	}

	return a, nil
}
