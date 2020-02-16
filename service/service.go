package service

import (
	user "github.com/frperezr/microservices-demo/src/users-api"
	"github.com/frperezr/microservices-demo/src/users-api/database"
)

// New ...
func New(store database.Store) *Users {
	return &Users{
		Store: store,
	}
}

// Users ...
type Users struct {
	Store database.Store
}

// GetByID ...
func (us *Users) GetByID(id string) (*user.User, error) {
	return us.Store.GetByID(id)
}

// GetByEmail ...
func (us *Users) GetByEmail(email string) (*user.User, error) {
	return us.Store.GetByEmail(email)
}

// Create ...
func (us *Users) Create(u *user.User) error {
	return us.Store.Create(u)
}

// Update ...
func (us *Users) Update(u *user.User) error {
	return us.Store.Update(u)
}

// Delete ...
func (us *Users) Delete(id string) error {
	return us.Store.Delete(id)
}
