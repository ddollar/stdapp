package stdapp

import (
	"context"
	"fmt"
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
	g := &GraphQL{
		app:    a,
		server: stdapi.New("api", a.opts.Name),
	}

	gopts := []graphqlws.Option{
		graphqlws.WithWriteTimeout(coalesce.Any(a.opts.WriteTimeout, 10*time.Second)),
	}

	for _, domain := range a.domains() {
		opts, err := pg.ParseURL(a.opts.Database)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		opts.PoolSize = 5

		domainCopy := domain

		opts.OnConnect = func(ctx context.Context, conn *pg.Conn) error {
			_, err := conn.Exec("SET search_path=?", domainCopy)
			return err
		}

		db := pg.Connect(opts)

		r, err := a.opts.Resolver(db, domain)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		h, err := stdgraph.NewHandler(r.Schema(), r, gopts...)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		g.server.Router.PathPrefix(fmt.Sprintf("%s/api/%s", a.opts.Prefix, domain)).Handler(a.WithMiddleware(h))

	}

	return g, nil
}
