package stdapp

import "github.com/go-pg/pg/v10/orm"

type ResolverFunc func(db orm.DB, domain string) (Resolver, error)

type Resolver interface {
	Mutation() any
	Query() any
	Schema() string
	Subscription() any
}
