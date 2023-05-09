package stdapp

type ResolverFunc func(db Database) (any, error)
