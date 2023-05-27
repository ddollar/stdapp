package stdapp

import (
	"net/http"
)

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func (a *App) WithMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(wrapMiddleware(h.ServeHTTP, a.opts.Middleware))
}

func wrapMiddleware(fn http.HandlerFunc, ms []Middleware) http.HandlerFunc {
	if len(ms) == 0 {
		return fn
	}

	return ms[0](wrapMiddleware(fn, ms[1:]))
}
