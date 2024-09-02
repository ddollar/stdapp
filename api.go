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

type API struct {
	app    *App
	server *stdapi.Server
}

func (a *App) api() (*API, error) {
	api := &API{
		app:    a,
		server: stdapi.New("api", a.opts.Name),
	}

	if err := api.handleGraphQL(a); err != nil {
		return nil, errors.Wrap(err)
	}

	if err := api.handleRouter(a); err != nil {
		return nil, errors.Wrap(err)
	}

	return api, nil
}

func (a *API) handleGraphQL(app *App) error {
	gopts := []graphqlws.Option{
		graphqlws.WithWriteTimeout(coalesce.Any(app.opts.WriteTimeout, 10*time.Second)),
	}

	for _, domain := range app.domains() {
		sdb := sql.OpenDB(pgdriver.NewConnector(
			pgdriver.WithDSN(app.opts.Database),
			pgdriver.WithConnParams(map[string]interface{}{
				"search_path": domain,
			}),
		))

		db := bun.NewDB(sdb, pgdialect.New())

		r, err := app.opts.Resolver(db, domain)
		if err != nil {
			return errors.Wrap(err)
		}

		h, err := stdgraph.NewHandler(r.Schema(), r, gopts...)
		if err != nil {
			return errors.Wrap(err)
		}

		a.server.Router.PathPrefix(fmt.Sprintf("%s/api/%s", app.opts.Prefix, domain)).Handler(app.WithMiddleware(h))
	}

	return nil
}

func (a *API) handleRouter(app *App) error {
	if app.opts.Router == nil {
		return nil
	}

	r := a.server.Subrouter(fmt.Sprintf("%s/api", app.opts.Prefix))

	if err := app.opts.Router(r); err != nil {
		return errors.Wrap(err)
	}

	return nil
}
