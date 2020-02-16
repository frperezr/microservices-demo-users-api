package users

import (
	"time"

	pb "github.com/frperezr/microservices-demo/pb"
)

// User is the main struct of the users api
type User struct {
	ID        string     `json:"id" db:"id"`
	Email     string     `json:"email" db:"email"`
	Name      string     `json:"name" db:"name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Password  string     `json:"password" db:"password"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// Service ...
type Service interface {
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(*User) error
	Update(*User) error
	Delete(id string) error
}

// ToProto ...
func (u *User) ToProto() *pb.User {
	return &pb.User{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		LastName:  u.LastName,
		Password:  u.Password,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}

// FromProto ...
func (u *User) FromProto(uu *pb.User) *User {
	return &User{
		ID:        uu.Id,
		Email:     uu.Email,
		Name:      uu.Name,
		LastName:  uu.LastName,
		Password:  uu.Password,
		CreatedAt: time.Unix(uu.CreatedAt, 0),
		UpdatedAt: time.Unix(uu.UpdatedAt, 0),
	}
}
