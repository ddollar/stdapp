package stdapp

import "github.com/uptrace/bun"

type ResolverFunc func(db *bun.DB, domain string) (Resolver, error)

type Resolver interface {
	Mutation() any
	Query() any
	Schema() string
	Subscription() any
}
