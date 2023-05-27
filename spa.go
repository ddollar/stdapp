package stdapp

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/ddollar/stdapi"
	"github.com/pkg/errors"
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

	prefix := fmt.Sprintf("%s/", a.opts.Prefix)

	s.server.Router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(http.FS(s))))

	return s, nil

}

func (s SPA) Open(name string) (fs.File, error) {
	if f, err := s.app.opts.Web.Open(name); err == nil {
		return f, nil
	} else {
		f, err := s.app.opts.Web.Open("index.html")
		if err != nil {
			return nil, errors.WithStack(err)
		}

		return f, nil
	}
}
