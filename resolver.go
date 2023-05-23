package stdapp

import "github.com/go-pg/pg/v10/orm"

type ResolverFunc func(db orm.DB) (any, error)
