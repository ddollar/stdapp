package stdapi

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"go.ddollar.dev/logger"
	"github.com/pkg/errors"
)

type RecoverFunc func(error)

type Server struct {
	*Router

	Check    HandlerFunc
	Hostname string
	Logger   *logger.Logger
	Recover  RecoverFunc
	Wrapper  func(h http.Handler) http.Handler

	middleware []Middleware
	server     http.Server
}

func (s *Server) HandleNotFound(fn HandlerFunc) {
	s.Router.HandleNotFound(fn)
}

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

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
