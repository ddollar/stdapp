package resolver

type Query struct {
	r *Resolver
}

func (r *Query) Ping() bool {
	return true
}
