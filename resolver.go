package stdapp

import "github.com/go-pg/pg/v10/orm"

type ResolverFunc func(db orm.DB) (Resolver, error)

type Resolver interface {
	Mutation() any
	Query() any
	Subscription() any
}
