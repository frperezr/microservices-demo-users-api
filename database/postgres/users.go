package postgres

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	user "github.com/frperezr/microservices-demo/src/users-api"
	"github.com/jmoiron/sqlx"
)

// UserStore ...
type UserStore struct {
	Store *sqlx.DB
}

// GetByID ...
func (us *UserStore) GetByID(id string) (*user.User, error) {
	if id == "" {
		return nil, errors.New("must provide a id")
	}

	query := squirrel.Select("*").From("users").Where("id = ? and deleted_at is null", id)

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := us.Store.QueryRowx(sql, args...)

	c := &user.User{}

	if err := row.StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// GetByEmail ...
func (us *UserStore) GetByEmail(email string) (*user.User, error) {
	if email == "" {
		return nil, errors.New("must provide a email")
	}

	query := squirrel.Select("*").From("users").Where("email = ? and deleted_at is null", email)

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := us.Store.QueryRowx(sql, args...)

	c := &user.User{}

	if err := row.StructScan(c); err != nil {
		return nil, err
	}

	return c, nil
}

// Create ...
func (us *UserStore) Create(u *user.User) error {
	if u.Email == "" {
		return errors.New("must provide a email")
	}

	sql, args, err := squirrel.
		Insert("users").
		Columns("email", "name", "last_name", "password").
		Values(strings.ToLower(u.Email), u.Name, u.LastName, u.Password).
		Suffix("returning *").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	row := us.Store.QueryRowx(sql, args...)
	if err := row.StructScan(u); err != nil {
		return err
	}

	return nil
}

// Update ...
func (us *UserStore) Update(u *user.User) error {
	if u.ID == "" {
		return errors.New("must provide a id")
	}

	query := squirrel.Update("users")

	if u.Email != "" {
		query = query.Set("email", strings.ToLower(u.Email))
	}

	if u.Name != "" {
		query = query.Set("name", u.Name)
	}

	if u.LastName != "" {
		query = query.Set("last_name", u.LastName)
	}

	if u.Password != "" {
		query = query.Set("password", u.Password)
	}

	sql, args, err := query.Where("id = ? and deleted_at is null", u.ID).Suffix("returning *").PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil
	}

	row := us.Store.QueryRowx(sql, args...)
	if err := row.StructScan(u); err != nil {
		return err
	}

	return nil
}

// Delete ...
func (us *UserStore) Delete(id string) error {
	if id == "" {
		return errors.New("must provide a id")
	}

	c := &user.User{}

	row := us.Store.QueryRowx("update users set deleted_at = $1 where id = $2", time.Now(), id)
	if err := row.StructScan(c); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}
