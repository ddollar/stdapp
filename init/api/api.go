package api

import (
	"example.org/stdapp/api/models"
	"github.com/go-pg/pg/v10/orm"
)

type API struct {
	models *models.Models
}

func New(db orm.DB) (any, error) {
	m, err := models.New(db)
	if err != nil {
		return nil, err
	}

	a := &API{
		models: m,
	}

	return a, nil
}

func (a *API) Ping() bool {
	return true
}
