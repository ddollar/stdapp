package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"time"

	"example.org/stdapp/api/resolver"
	"github.com/ddollar/errors" //go:embed db/migrate/*.sql
	"github.com/ddollar/stdapp"
)

var migrations embed.FS

//go:embed all:web/dist
var web embed.FS

func main() {
	a, err := app()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	os.Exit(a.Run(os.Args[1:]))
}

func app() (*stdapp.App, error) {
	sweb, err := fs.Sub(fs.FS(web), "web/dist")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	opts := stdapp.Options{
		Database:     os.Getenv("DATABASE_URL"),
		Migrations:   migrations,
		Name:         "stdapp-init",
		Resolver:     resolver.New,
		Web:          sweb,
		WriteTimeout: 5 * time.Minute,
	}

	a, err := stdapp.New(opts)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return a, nil
}
