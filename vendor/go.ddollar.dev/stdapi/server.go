package stdapi

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"go.ddollar.dev/logger"
	"github.com/pkg/errors"
)

// RecoverFunc is called when a panic occurs during request handling.
// It receives the recovered error and can be used for logging or alerting.
type RecoverFunc func(error)

// Server represents an HTTP/HTTPS server with routing capabilities.
//
// Server embeds Router and provides additional server-level configuration
// including health checks, panic recovery, and middleware wrapping.
type Server struct {
	*Router

	// Check is an optional custom health check handler.
	// If set, it will be called by the /check endpoint.
	Check HandlerFunc

	// Hostname is the server hostname used for TLS certificates.
	Hostname string

	// Logger is the structured logger for this server.
	Logger *logger.Logger

	// Recover is an optional panic recovery handler.
	// If set, panics during request handling will be passed to this function.
	Recover RecoverFunc

	// Wrapper is an optional middleware that wraps the entire server handler.
	Wrapper func(h http.Handler) http.Handler

	middleware []Middleware
	server     http.Server
}

// HandleNotFound sets a custom handler for 404 Not Found responses.
func (s *Server) HandleNotFound(fn HandlerFunc) {
	s.Router.HandleNotFound(fn)
}

// Listen starts the HTTP server on the specified protocol and address.
//
// Supported protocols:
//   - "http" - Plain HTTP
//   - "https", "tls" - HTTPS with auto-generated self-signed certificate
//   - "h2" - HTTP/2 with TLS
//
// The address should be in the format "host:port" (e.g., ":8080" or "localhost:443").
//
// For TLS protocols, a self-signed certificate is automatically generated using
// the server's Hostname field.
//
// This method blocks until the server is shut down or encounters an error.
func (s *Server) Listen(proto, addr string) error {
	s.Logger.At("listen").Logf("hostname=%q proto=%q addr=%q", s.Hostname, proto, addr)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return errors.WithStack(err)
	}

	switch proto {
	case "h2", "https", "tls":
		config := &tls.Config{
			MinVersion: tls.VersionTLS12,
		}

		if proto == "h2" {
			config.NextProtos = []string{"h2"}
		}

		cert, err := generateSelfSignedCertificate(s.Hostname)
		if err != nil {
			return errors.WithStack(err)
		}

		config.Certificates = append(config.Certificates, cert)

		l = tls.NewListener(l, config)
	}

	var h http.Handler

	if s.Wrapper != nil {
		h = s.Wrapper(s)
	} else {
		h = s
	}

	s.server = http.Server{Handler: h}

	return s.server.Serve(l)
}

// Shutdown gracefully shuts down the server without interrupting active connections.
//
// The context determines the maximum time to wait for active connections to close.
// Returns an error if the context deadline is exceeded or another shutdown error occurs.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
