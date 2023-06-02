package main

import (
	"fmt"
	"os"

	"github.com/ddollar/stdapp"
)

type compose struct {
	Name     string
	Services map[string]struct {
		Environment map[string]string
	}
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
	// if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
	// 	return nil, fmt.Errorf("must run in a directory with docker-compose.yml")
	// }

	// data, err := exec.Command("docker", "compose", "config", "--format=json").CombinedOutput()
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// var c compose

	// if err := json.Unmarshal(data, &c); err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// svc := coalesce.Any(c.Services["api"], c.Services["web"])
	// env := svc.Environment
	// dburl := coalesce.Any(env["POSTGRES_URL"], env["DATABASE_URL"])

	// opts := &stdapp.Options{
	// 	Compose:  true,
	// 	Database: dburl,
	// 	Name:     c.Name,
	// }

	opts := &stdapp.Options{
		Compose: true,
		Name:    "foo",
	}

	return opts, nil
}
