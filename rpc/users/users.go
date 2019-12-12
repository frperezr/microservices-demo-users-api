package users

import (
	"fmt"
	"log"

	pb "github.com/frperezr/noken-test/pb"
	"github.com/frperezr/noken-test/src/users-api"
	"github.com/frperezr/noken-test/src/users-api/database"
	"github.com/frperezr/noken-test/src/users-api/service"
	"golang.org/x/net/context"
)

var _ pb.UserServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	userSvc users.Service
}

// New ...
func New(store database.Store) *Service {
	return &Service{
		userSvc: service.New(store),
	}
}

// GetByID ...
func (us *Service) GetByID(ctx context.Context, gr *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	id := gr.GetId()
	log.Println(fmt.Sprintf("[User Service][GetById][Request] id = %v", id))

	if id == "" {
		log.Println("[User Service][GetById][Error] must provide a id")
		return &pb.GetUserByIDResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "must provide a id",
			},
		}, nil
	}

	user, err := us.userSvc.GetByID(id)
	if err != nil {
		log.Println(fmt.Sprintf("[User Service][GetById][Error] %v", err.Error()))
		if err.Error() == "sql: no rows in result set" {
			return &pb.GetUserByIDResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    404,
					Message: fmt.Sprintf("user with id %v not found", id),
				},
			}, nil
		}

		return &pb.GetUserByIDResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.GetUserByIDResponse{
		Data:  user.ToProto(),
		Error: nil,
	}

	log.Println(fmt.Sprintf("[User Service][GetById][Response] %v", res))
	return res, nil
}

// GetByEmail ...
func (us *Service) GetByEmail(ctx context.Context, gr *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	email := gr.GetEmail()
	log.Println(fmt.Sprintf("[User Service][GetByEmail][Request] email = %v", email))

	if email == "" {
		log.Println("[User Service][GetByEmail][Error] must provide a email")
		return &pb.GetUserByEmailResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "must provide a email",
			},
		}, nil
	}

	user, err := us.userSvc.GetByEmail(email)
	if err != nil {
		log.Println(fmt.Sprintf("[User Service][GetByEmail][Error] %v", err.Error()))
		if err.Error() == "sql: no rows in result set" {
			return &pb.GetUserByEmailResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    404,
					Message: fmt.Sprintf("user with email %v not found", email),
				},
			}, nil
		}

		return &pb.GetUserByEmailResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.GetUserByEmailResponse{
		Data:  user.ToProto(),
		Error: nil,
	}

	log.Println(fmt.Sprintf("[User Service][GetByEmail][Response] %v", res))
	return res, nil
}

// Create ...
func (us *Service) Create(ctx context.Context, gr *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	email := gr.GetData().GetEmail()
	log.Println(fmt.Sprintf("[User Service][Create][Request] email = %v", email))

	if email == "" {
		log.Println("[User Service][Create][Error] email param is empty")
		return &pb.CreateUserResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "email param is empty",
			},
		}, nil
	}

	_, err := us.userSvc.GetByEmail(email)
	if err != nil {

		name := gr.GetData().GetName()
		if name == "" {
			log.Println(fmt.Sprintf("[User Service][Create][Error] %v", "name param is empty"))
			return &pb.CreateUserResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    400,
					Message: "name param is empty",
				},
			}, nil
		}

		lastName := gr.GetData().GetLastName()
		if lastName == "" {
			log.Println(fmt.Sprintf("[User Service][Create][Error] %v", "last_name param is empty"))
			return &pb.CreateUserResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    400,
					Message: "last_name param is empty",
				},
			}, nil
		}

		password := gr.GetData().GetPassword()
		if password == "" {
			log.Println(fmt.Sprintf("[User Service][Create][Error] %v", "password param is empty"))
			return &pb.CreateUserResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    400,
					Message: "password param is empty",
				},
			}, nil
		}

		user := &users.User{
			Email:    email,
			Name:     name,
			LastName: lastName,
			Password: password,
		}

		if err := us.userSvc.Create(user); err != nil {
			log.Println(fmt.Sprintf("[User Service][Create][Error] %v", err.Error()))
			return &pb.CreateUserResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    500,
					Message: err.Error(),
				},
			}, nil
		}

		res := &pb.CreateUserResponse{
			Data:  user.ToProto(),
			Error: nil,
		}
		log.Println(fmt.Sprintf("[User Service][Create][Response] %v", res))
		return res, nil

	}

	res := &pb.CreateUserResponse{
		Data: nil,
		Error: &pb.Error{
			Code:    500,
			Message: "user already registered",
		},
	}
	log.Println(fmt.Sprintf("[User Service][Create][Response] %v", res))
	return res, nil
}

// Update ...
func (us *Service) Update(ctx context.Context, gr *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	id := gr.GetData().GetId()
	log.Println(fmt.Sprintf("[User Service][Update][Request] id = %v", id))

	if id == "" {
		log.Println("[User Service][Update][Error] id param is empty")
		return &pb.UpdateUserResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    400,
				Message: "id param is empty",
			},
		}, nil
	}

	user, err := us.userSvc.GetByID(id)
	if err != nil {
		log.Println(fmt.Sprintf("[User Service][Update][Error] err = %v", err.Error()))

		if err.Error() == "sql: no rows in result set" {
			return &pb.UpdateUserResponse{
				Data: nil,
				Error: &pb.Error{
					Code:    404,
					Message: fmt.Sprintf("user with id = %v not found", id),
				},
			}, nil
		}

		return &pb.UpdateUserResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	email := gr.GetData().GetEmail()
	if email != "" {
		user.Email = email
	}

	name := gr.GetData().GetName()
	if name != "" {
		user.Name = name
	}

	lastName := gr.GetData().GetLastName()
	if lastName != "" {
		user.LastName = lastName
	}

	password := gr.GetData().GetPassword()
	if password != "" {
		user.Password = password
	}

	if err := us.userSvc.Update(user); err != nil {
		log.Println(fmt.Sprintf("[User Service][Update][Error] err = %v", err.Error()))
		return &pb.UpdateUserResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.UpdateUserResponse{
		Data:  user.ToProto(),
		Error: nil,
	}

	log.Println(fmt.Sprintf("[User Service][Update][Response] %v", res))
	return res, nil
}

// Delete ...
func (us *Service) Delete(ctx context.Context, gr *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	id := gr.GetUserId()
	log.Println(fmt.Sprintf("[User Service][Delete][Request] id = %v", id))

	user, err := us.userSvc.GetByID(id)
	if err != nil {
		log.Println(fmt.Sprintf("[User Service][Delete][Error] err = %v", err.Error()))
		return &pb.DeleteUserResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	if err := us.userSvc.Delete(user.ID); err != nil {
		log.Println(fmt.Sprintf("[User Service][Delete][Error] err = %v", err.Error()))
		return &pb.DeleteUserResponse{
			Data: nil,
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	res := &pb.DeleteUserResponse{
		Data:  user.ToProto(),
		Error: nil,
	}

	log.Println(fmt.Sprintf("[User Service][Delete][Response] %v", res))
	return res, nil
}
