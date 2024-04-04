package stdapp

import (
	"io/fs"
	"net/http"

	"github.com/ddollar/errors"
	"github.com/ddollar/stdapi"
)

type SPA struct {
	app    *App
	server *stdapi.Server
}

func (a *App) spa() (*SPA, error) {
	s := &SPA{
		app:    a,
		server: stdapi.New(a.opts.Name, a.opts.Name),
	}

	h := http.StripPrefix(a.opts.Prefix, http.FileServer(http.FS(s)))

	s.server.Router.PathPrefix(a.opts.Prefix).Handler(a.WithMiddleware(h))

	return s, nil

}

func (s SPA) Open(name string) (fs.File, error) {
	if f, err := s.app.opts.Web.Open(name); err == nil {
		return f, nil
	} else {
		f, err := s.app.opts.Web.Open("index.html")
		if err != nil {
			return nil, errors.Wrap(err)
		}

		return f, nil
	}
}
