package database

import (
	user "github.com/frperezr/microservices-demo/src/users-api"
	"github.com/frperezr/microservices-demo/src/users-api/database/postgres"
	"github.com/jmoiron/sqlx"
)

// Store ...
type Store interface {
	GetByID(id string) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Create(*user.User) error
	Update(*user.User) error
	Delete(id string) error
}

// NewPostgres ...
func NewPostgres(dsn string) (Store, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &postgres.UserStore{
		Store: db,
	}, nil
}
