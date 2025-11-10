//go:build tools

package stdapp

import (
	_ "github.com/cespare/reflex"
	_ "go.ddollar.dev/stdapp"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
