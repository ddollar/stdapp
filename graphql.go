package stdapp

import (
	"github.com/ddollar/stdapi"
	"github.com/ddollar/stdgraph"
	"github.com/pkg/errors"
)

type GraphQL struct {
	app    *App
	server *stdapi.Server
}

func (a *App) graphQL() (*GraphQL, error) {
	g := &GraphQL{
		app:    a,
		server: stdapi.New(a.name, a.name),
	}

	h, err := stdgraph.NewHandler(a.schema, a.resolver)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	g.server.Router.PathPrefix("/api/graph").Handler(h)

	return g, nil
}
