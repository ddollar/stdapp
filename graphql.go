package stdapp

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ddollar/coalesce"
	"github.com/ddollar/errors"
	"github.com/ddollar/graphql-transport-ws/graphqlws"
	"github.com/ddollar/stdapi"
	"github.com/ddollar/stdgraph"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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
		sdb := sql.OpenDB(pgdriver.NewConnector(
			pgdriver.WithDSN(a.opts.Database),
			pgdriver.WithConnParams(map[string]interface{}{
				"search_path": domain,
			}),
		))

		db := bun.NewDB(sdb, pgdialect.New())

		r, err := a.opts.Resolver(db, domain)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		h, err := stdgraph.NewHandler(r.Schema(), r, gopts...)
		if err != nil {
			return nil, errors.Wrap(err)
		}

		g.server.Router.PathPrefix(fmt.Sprintf("%s/api/%s", a.opts.Prefix, domain)).Handler(a.WithMiddleware(h))

	}

	return g, nil
}
