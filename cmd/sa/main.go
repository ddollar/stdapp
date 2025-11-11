package main

import (
	"fmt"
	"os"
	"path/filepath"

	"go.ddollar.dev/errors"
	"go.ddollar.dev/stdapp"
	"gopkg.in/yaml.v3"
)

type kip struct {
	Resources map[string]struct {
		Type string `yaml:"type"`
	} `yaml:"resources"`
	Services map[string]struct {
		Environment map[string]string `yaml:"environment"`
	} `yaml:"services"`
}

func main() {
	opts, err := options()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	a, err := stdapp.New(*opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	os.Exit(a.Run(os.Args[1:]))
}

func options() (*stdapp.Options, error) {
	opts := &stdapp.Options{
		Compose: true,
	}

	if _, err := os.Stat("kip.yml"); os.IsNotExist(err) {
		return opts, nil
	}

	data, err := os.ReadFile("kip.yml")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var k kip

	if err := yaml.Unmarshal(data, &k); err != nil {
		return nil, errors.Wrap(err)
	}

	// Use directory name as app name
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err)
	}
	opts.Name = filepath.Base(cwd)

	// Check if postgres resource exists to set DATABASE_URL
	if _, ok := k.Resources["postgres"]; ok {
		// If DATABASE_URL is not already set in environment, kip will inject it
		if dburl := os.Getenv("DATABASE_URL"); dburl != "" {
			opts.Database = dburl
		}
	}

	// Check services for explicit DATABASE_URL
	if api, ok := k.Services["api"]; ok {
		if dburl := api.Environment["DATABASE_URL"]; dburl != "" {
			opts.Database = dburl
		}
	}
	if web, ok := k.Services["web"]; ok {
		if dburl := web.Environment["DATABASE_URL"]; dburl != "" {
			opts.Database = dburl
		}
	}

	return opts, nil
}
