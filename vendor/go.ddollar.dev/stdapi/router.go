package stdapi

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// HandlerFunc is the signature for HTTP request handlers in stdapi.
//
// Unlike standard http.HandlerFunc, stdapi handlers receive a Context object
// that wraps the request and response, and they return an error instead of
// writing errors directly to the response.
//
// If a handler returns an error that implements the Error interface, the error's
// Code() will be used as the HTTP status code. Otherwise, a 500 Internal Server
// Error is returned.
type HandlerFunc func(c *Context) error

// Middleware wraps a HandlerFunc to add pre- or post-processing logic.
//
// Middleware can modify the request/response, log information, check authentication,
// or short-circuit the handler chain by returning early.
type Middleware func(fn HandlerFunc) HandlerFunc

// Router handles HTTP routing and middleware management.
//
// Router wraps gorilla/mux.Router and provides the stdapi-style handler interface.
// Routers can be nested using Subrouter, and child routers inherit parent middleware.
type Router struct {
	*mux.Router

	// Middleware contains the middleware stack for this router.
	Middleware []Middleware

	// Parent is the parent router if this is a subrouter, nil otherwise.
	Parent *Router

	// Server is the server this router belongs to.
	Server *Server
}

// Route represents a registered route and provides access to underlying mux.Route
// for additional configuration (e.g., host matching, schemes, queries).
type Route struct {
	*mux.Route

	Router *Router
}

// MatcherFunc creates a subrouter that only matches requests satisfying the given matcher function.
//
// The matcher function is called for each request and should return true if the request
// should be handled by this subrouter. This is useful for custom routing logic like
// header-based routing or complex conditional routing.
func (rt *Router) MatcherFunc(fn mux.MatcherFunc) *Router {
	return &Router{
		Parent: rt,
		Router: rt.Router.MatcherFunc(fn).Subrouter(),
		Server: rt.Server,
	}
}

// HandleNotFound sets a custom handler for requests that don't match any routes.
func (rt *Router) HandleNotFound(fn HandlerFunc) {
	rt.Router.NotFoundHandler = rt.http(fn)
}

// Redirect registers a route that redirects to the target URL with the given status code.
//
// Common redirect codes:
//   - 301 - Permanent redirect
//   - 302 - Temporary redirect
//   - 307 - Temporary redirect (preserves method)
//   - 308 - Permanent redirect (preserves method)
func (rt *Router) Redirect(method, path string, code int, target string) {
	rt.Handle(path, Redirect(code, target)).Methods(method)
}

// Route registers a handler for the given HTTP method and path pattern.
//
// Supported methods:
//   - "GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS" - Standard HTTP methods
//   - "SOCKET" - WebSocket endpoint (automatically handles upgrade handshake)
//   - "ANY" - Matches any HTTP method
//
// Path patterns support gorilla/mux syntax including:
//   - Static paths: "/users"
//   - Path variables: "/users/{id}"
//   - Regular expressions: "/users/{id:[0-9]+}"
//   - Path prefixes via Subrouter
//
// The returned Route can be further configured using mux.Route methods.
func (rt *Router) Route(method, path string, fn HandlerFunc) Route {
	switch method {
	case "SOCKET":
		return Route{
			Route:  rt.Handle(path, rt.websocket(fn)).Methods("GET").Headers("Upgrade", "websocket"),
			Router: rt,
		}
	case "ANY":
		return Route{
			Route:  rt.Handle(path, rt.http(fn)),
			Router: rt,
		}
	default:
		return Route{
			Route:  rt.Handle(path, rt.http(fn)).Methods(method),
			Router: rt,
		}
	}
}

// Static serves static files from the given filesystem at the specified path prefix.
//
// Example:
//
//	s.Static("/assets", http.Dir("./public"))
//
// This would serve files from ./public at URLs like /assets/style.css
func (rt *Router) Static(path string, files FileSystem) Route {
	prefix := fmt.Sprintf("%s/", path)

	return Route{
		Route:  rt.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(files))),
		Router: rt,
	}
}

// Subrouter creates a new router that handles paths with the given prefix.
//
// The subrouter inherits all middleware from its parent routers, executing them
// in order from outermost to innermost.
//
// Example:
//
//	api := s.Subrouter("/api/v1")
//	api.Route("GET", "/users", listUsers)  // Handles /api/v1/users
func (rt *Router) Subrouter(prefix string) *Router {
	return &Router{
		Parent: rt,
		Router: rt.PathPrefix(prefix).Subrouter(),
		Server: rt.Server,
	}
}

// SubrouterFunc creates a subrouter and immediately calls the given function with it.
//
// This is a convenience method for organizing route registration:
//
//	s.SubrouterFunc("/api", func(api *Router) {
//		api.Route("GET", "/users", listUsers)
//		api.Route("POST", "/users", createUser)
//	})
func (rt *Router) SubrouterFunc(prefix string, fn func(*Router)) {
	fn(rt.Subrouter(prefix))
}

// Use adds middleware to this router's middleware stack.
//
// Middleware is executed in the order it's added, with parent router middleware
// executing before child router middleware.
func (rt *Router) Use(mw Middleware) {
	rt.Middleware = append(rt.Middleware, mw)
}

// UseHandlerFunc adds standard http.HandlerFunc middleware to this router.
//
// This allows using standard Go HTTP middleware with stdapi routers.
// The HandlerFunc receives the raw http.ResponseWriter and *http.Request.
func (rt *Router) UseHandlerFunc(fn http.HandlerFunc) {
	rt.Middleware = append(rt.Middleware, func(gn HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			fn(c.response, c.request)
			return gn(c)
		}
	})
}

func (rt *Router) context(name string, w http.ResponseWriter, r *http.Request, conn *websocket.Conn) (*Context, error) {
	id, err := generateId(12)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	w.Header().Add("Request-Id", id)

	c := NewContext(w, r)

	c.context = context.WithValue(r.Context(), "request.id", id)
	c.id = id
	c.logger = rt.Server.Logger.Prepend("id=%s", id).At(name)
	c.name = name
	c.ws = conn

	return c, nil
}

func (rt *Router) handle(fn HandlerFunc, c *Context) error {
	defer func() {
		if rt.Server.Recover == nil {
			return
		}
		switch t := recover().(type) {
		case error:
			rt.Server.Recover(t)
		case string:
			rt.Server.Recover(fmt.Errorf(t))
		case fmt.Stringer:
			rt.Server.Recover(fmt.Errorf(t.String()))
		case nil:
			return
		default:
			panic(t)
		}
	}()

	c.logger = c.logger.Append("method=%q path=%q", c.request.Method, c.request.URL.Path).Start()

	// rw := &responseWriter{ResponseWriter: c.response, code: 200}
	// c.response = rw

	mw := []Middleware{}

	p := rt.Parent

	for {
		if p == nil {
			break
		}

		mw = append(p.Middleware, mw...)

		p = p.Parent
	}

	mw = append(mw, rt.Middleware...)

	fnmw := rt.wrap(fn, mw...)

	errr := fnmw(c) // non-standard error name to avoid wrapping

	if ne, ok := errr.(*net.OpError); ok {
		c.logger.Logf("state=closed error=%q", ne.Err)
		return nil
	}

	code := c.response.Code()

	if code == 0 {
		code = 200
	}

	c.logger.Logf("response=%d", code)

	return errr
}

func (rt *Router) http(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := rt.context(functionName(fn), w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if err := rt.handle(fn, c); err != nil {
			switch t := err.(type) {
			case Error:
				c.logger.Append("code=%d", t.Code()).Error(err)
				http.Error(c.response, t.Error(), t.Code())
			case error:
				c.logger.Error(err)
				http.Error(c.response, t.Error(), http.StatusInternalServerError)
			}
		}
	}
}

var upgrader = websocket.Upgrader{ReadBufferSize: 10 * 1024, WriteBufferSize: 10 * 1024}

func (rt *Router) websocket(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("ERROR: %s\n", err.Error())))
			return
		}

		// empty binary message signals EOF
		defer conn.WriteMessage(websocket.BinaryMessage, []byte{})
		// defer conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""))

		c, err := rt.context(functionName(fn), w, r, conn)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ERROR: %s\n", err.Error())))
			return
		}

		if err := rt.handle(fn, c); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("ERROR: %s\n", err.Error())))
			return
		}
	}
}

func (rt *Router) wrap(fn HandlerFunc, m ...Middleware) HandlerFunc {
	if len(m) == 0 {
		return fn
	}

	return m[0](rt.wrap(fn, m[1:len(m)]...))
}
