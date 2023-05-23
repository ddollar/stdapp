package stdapp

import (
	"fmt"
	"strings"

	"github.com/ddollar/logger"
)

func New(opts Options) (*App, error) {
	a := &App{
		compose:    opts.Compose,
		database:   opts.Database,
		logger:     logger.New(fmt.Sprintf("ns=%s", opts.Name)),
		migrations: opts.Migrations,
		name:       opts.Name,
		schema:     opts.Schema,
	}

	return a, nil
}

func parseExtensions(flag string) []string {
	return strings.Split(flag, ",")
}
