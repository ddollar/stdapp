// Package stdapi provides a lightweight, opinionated HTTP API framework for Go.
//
// stdapi is built on top of gorilla/mux and provides a cleaner, more ergonomic API
// for building web applications and HTTP APIs. It wraps standard Go HTTP handlers
// with a custom Context object and structured error handling.
//
// Key features:
//   - Context-based request/response handling
//   - Structured error handling with HTTP status codes
//   - Built-in middleware support
//   - WebSocket support alongside HTTP endpoints
//   - Session management with flash messages
//   - Template rendering with layout support
//   - Automatic request parameter unmarshaling
//   - TLS/HTTP2 support with auto-generated certificates
//   - Request ID generation and structured logging
//
// Example usage:
//
//	s := stdapi.New("myapp", "localhost")
//	s.Route("GET", "/users", listUsers)
//	s.Route("POST", "/users", createUser)
//	s.Listen("https", ":443")
//
//	func listUsers(c *stdapi.Context) error {
//		users := []User{...}
//		return c.RenderJSON(users)
//	}
package stdapi

import (
	"fmt"
	"net/http"

	"go.ddollar.dev/logger"
	"github.com/gorilla/mux"
)

// New creates a new Server with the given namespace and hostname.
//
// The namespace is used for structured logging to identify the application.
// The hostname is used for TLS certificate generation when using secure protocols.
//
// A default health check endpoint is automatically registered at /check.
func New(ns, hostname string) *Server {
	logger := logger.New(fmt.Sprintf("ns=%s", ns))

	server := &Server{
		Hostname: hostname,
		Logger:   logger,
	}

	server.Router = &Router{
		Parent: nil,
		Router: mux.NewRouter(),
		Server: server,
	}

	server.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		id, _ := generateId(12)
		logger.Logf("id=%s route=unknown code=404 method=%q path=%q", id, r.Method, r.URL.Path)
	})

	server.Router.HandleFunc("/check", server.check)

	return server
}

func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	if s.Check != nil {
		if err := s.Check(NewContext(w, r)); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	fmt.Fprintf(w, "ok")
}
