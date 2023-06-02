package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/netip"
	"os"
	"time"

	"example.org/stdapp/api"
	"github.com/ddollar/stdapp"
	"github.com/pkg/errors"
)

//go:embed db/migrate/*.sql
var migrations embed.FS

//go:embed api/schema.graphql
var schema string

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
		return nil, errors.WithStack(err)
	}

	opts := stdapp.Options{
		Database:     os.Getenv("DATABASE_URL"),
		Middleware:   []stdapp.Middleware{ensureAllowedNetwork},
		Migrations:   migrations,
		Name:         "stdapp-init",
		Resolver:     api.New,
		Schema:       schema,
		Web:          sweb,
		WriteTimeout: 5 * time.Minute,
	}

	a, err := stdapp.New(opts)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return a, nil
}

var allowedNetworks = []netip.Prefix{
	netip.MustParsePrefix("10.1.48.0/22"),
}

func ensureAllowedNetwork(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		remote := r.Header.Get("X-Forwarded-For")

		if os.Getenv("DEVELOPMENT") == "true" {
			fmt.Printf("ns=network at=ensure remote=%q development=true\n", remote)
			next(w, r)
			return
		}

		addr, err := netip.ParseAddr(remote)
		if err != nil {
			fmt.Printf("ns=network at=ensure remote=%q error=%q\n", remote, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, net := range allowedNetworks {
			if net.Contains(addr) {
				fmt.Printf("ns=network at=ensure remote=%q allowed=true net=%q\n", remote, net)
				next(w, r)
				return
			}
		}

		fmt.Printf("ns=network at=ensure remote=%q allowed=false\n", remote)
		http.Error(w, "not authorized", http.StatusUnauthorized)
	}
}
