package stdapp

import (
	"strings"

	"github.com/ddollar/logger"
)

func New(opts Options) (*App, error) {
	a := &App{
		opts:   opts,
		logger: logger.New("ns=stdapp"),
	}

	return a, nil
}

func parseExtensions(flag string) []string {
	return strings.Split(flag, ",")
}
