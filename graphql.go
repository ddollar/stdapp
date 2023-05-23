package stdapp

import (
	"github.com/ddollar/stdapi"
	"github.com/ddollar/stdgraph"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

type GraphQL struct {
	app    *App
	server *stdapi.Server
}

func (a *App) graphQL() (*GraphQL, error) {
	opts, err := pg.ParseURL(a.database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	opts.PoolSize = 5

	db := pg.Connect(opts)

	r, err := a.resolver(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	g := &GraphQL{
		app:    a,
		server: stdapi.New(a.name, a.name),
	}

	h, err := stdgraph.NewHandler(a.schema, r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	g.server.Router.PathPrefix("/api/graph").Handler(h)

	return g, nil
}
