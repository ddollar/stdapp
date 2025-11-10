package stdapp

import (
	"strings"

	"go.ddollar.dev/logger"
)

func New(opts Options) (*App, error) {
	a := &App{
		opts:   opts,
		logger: logger.New("ns=stdapp"),
	}

	return a, nil
}

func parseExtensions(flag string) []string {
	if strings.TrimSpace(flag) == "" {
		return []string{}
	}

	return strings.Split(flag, ",")
}
