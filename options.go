package stdapp

import (
	"io/fs"
)

type Options struct {
	Compose    bool
	Database   string
	Middleware []Middleware
	Migrations fs.FS
	Name       string
	Resolver   ResolverFunc
	Schema     string
	Web        fs.FS
}
