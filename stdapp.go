package stdapp

import (
	"fmt"
	"strings"

	"github.com/ddollar/logger"
)

func New(opts Options) (*App, error) {
	a := &App{
		opts:   opts,
		logger: logger.New(fmt.Sprintf("ns=%s", opts.Name)),
		// compose:    opts.Compose,
		// database:   opts.Database,
		// middleware: opts.Middleware,
		// migrations: opts.Migrations,
		// name:       opts.Name,
		// resolver:   opts.Resolver,
		// schema:     opts.Schema,
		// web:        opts.Web,
	}

	return a, nil
}

func parseExtensions(flag string) []string {
	return strings.Split(flag, ",")
}
