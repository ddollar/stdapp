package stdapp

import (
	"time"

	"github.com/ddollar/coalesce"
	"github.com/ddollar/graphql-transport-ws/graphqlws"
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
	opts, err := pg.ParseURL(a.opts.Database)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	opts.PoolSize = 5

	db := pg.Connect(opts)

	r, err := a.opts.Resolver(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	g := &GraphQL{
		app:    a,
		server: stdapi.New(a.opts.Name, a.opts.Name),
	}

	gopts := []graphqlws.Option{
		graphqlws.WithWriteTimeout(coalesce.Any(a.opts.WriteTimeout, 10*time.Second)),
	}

	h, err := stdgraph.NewHandler(a.opts.Schema, r, gopts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	g.server.Router.PathPrefix("/api/graph").Handler(a.WithMiddleware(h))

	return g, nil
}
