package stdapp

import (
	"github.com/ddollar/stdapi"
	"github.com/ddollar/stdgraph"
	"golang.org/x/xerrors"
)

type Api struct {
	server *stdapi.Server
}

func NewApi(schema string, resolver any) (*Api, error) {
	a := &Api{
		server: stdapi.New("stdapp", "stdapp"),
	}

	g, err := stdgraph.NewHandler(schema, resolver)
	if err != nil {
		return nil, err
	}

	a.server.Router.PathPrefix("/api/graph").Handler(g)

	return a, nil
}

func (a *Api) ListenAndServe(addr string) error {
	if err := a.server.Listen("https", addr); err != nil {
		return xerrors.Errorf("listen: %w", err)
	}

	return nil
}
