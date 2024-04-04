package resolver

import (
	_ "embed"

	"example.org/stdapp/api/models"
	"github.com/ddollar/stdapp"
	"github.com/go-pg/pg/v10/orm"
)

//go:embed schema.graphql
var schema string

type Resolver struct {
	models *models.Models
}

func New(db orm.DB, domain string) (stdapp.Resolver, error) {
	m, err := models.New(db)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r := &Resolver{
		models: m,
	}

	return r, nil
}

func (r *Resolver) Mutation() any {
	return &Mutation{r: r}
}

func (r *Resolver) Query() any {
	return &Query{r: r}
}

func (r *Resolver) Schema() string {
	return schema
}

func (r *Resolver) Subscription() any {
	return &Subscription{r: r}
}

type Wrapper[M any] struct {
	r    *Resolver
	item M
}

func wrap[M any](r *Resolver, item M) Wrapper[M] {
	return Wrapper[M]{r: r, item: item}
}
