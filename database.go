package stdapp

import (
	"net/url"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/pkg/errors"
)

type Database interface {
	Ping() error
}

func database(url_ string) (Database, error) {
	u, err := url.Parse(url_)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch u.Scheme {
	case "postgres":
		return pgInitialize(url_)
	default:
		return nil, errors.Errorf("unknown database scheme: %s", u.Scheme)
	}
}

type Postgres struct {
	db orm.DB
}

func pgInitialize(url_ string) (*Postgres, error) {
	opts, err := pg.ParseURL(url_)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	opts.PoolSize = 5

	p := &Postgres{
		db: pg.Connect(opts),
	}

	return p, nil
}

func (p *Postgres) Ping() error {
	_, err := p.db.Exec("SELECT 1")
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
