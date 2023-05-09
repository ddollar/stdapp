package stdapp

import "io/fs"

type Options struct {
	Database   string
	Migrations fs.FS
	Name       string
	Resolver   ResolverFunc
	Schema     string
}
