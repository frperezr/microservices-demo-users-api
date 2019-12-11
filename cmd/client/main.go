package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	pb "github.com/frperezr/noken-test/pb"
	"github.com/frperezr/noken-test/src/users-api"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	usersHost := os.Getenv("USERS_HOST")
	if usersHost == "" {
		fmt.Print(`{"error": "missing env USERS_HOST"}`)
		os.Exit(1)
	}

	usersPort := os.Getenv("USERS_PORT")
	if usersPort == "" {
		fmt.Print(`{"error": "missing env USERS_PORT"}`)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%v:%v", usersHost, usersPort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	subcmd := flag.Arg(0)
	var result string
	switch subcmd {
	case "getById":
		result, err = GetByID(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
			os.Exit(1)
		}
	case "getByEmail":
		result, err = GetByEmail(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
			os.Exit(1)
		}
	case "create":
		result, err = Create(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
			os.Exit(1)
		}
	case "delete":
		result, err = Delete(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
			os.Exit(1)
		}
	case "update":
		result, err = Update(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
			os.Exit(1)
		}
	default:
		fmt.Print(`{"error": "invalid command"}`)
		os.Exit(1)
	}

	fmt.Print(result)
	os.Exit(0)
}

// GetByID returns a user by ID.
func GetByID(us pb.UserServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing id param")
	}

	jsonStr := args[0]
	id := struct {
		ID string `json:"id"`
	}{}

	err := json.Unmarshal([]byte(jsonStr), &id)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := us.GetByID(context.Background(), &pb.GetUserByIDRequest{
		Id: id.ID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", fmt.Errorf("user with id = %v not found", id.ID)
		}

		return "", fmt.Errorf(err.Error())
	}

	if res.Error != nil {
		if strings.Contains(res.Error.Message, "sql: no rows in result set") {
			return "", fmt.Errorf("user with id = %v not found", id.ID)
		}

		return "", fmt.Errorf(err.Error())
	}

	json, err := json.Marshal(res.Data)
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}

// GetByEmail returns a user by email
func GetByEmail(us pb.UserServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing email param")
	}

	jsonStr := args[0]
	email := struct {
		Email string `json:"email"`
	}{}

	err := json.Unmarshal([]byte(jsonStr), &email)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := us.GetByEmail(context.Background(), &pb.GetUserByEmailRequest{
		Email: email.Email,
	})

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", fmt.Errorf("user with email = %v not found", email.Email)
		}

		return "", errors.New(err.Error())
	}

	json, err := json.Marshal(res.Data)
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}

// Create makes a new user
func Create(us pb.UserServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing user param")
	}

	jsonStr := args[0]
	user := users.User{}

	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := us.Create(context.Background(), &pb.CreateUserRequest{
		Data: user.ToProto(),
	})

	if err != nil {
		return "", errors.New(err.Error())
	}

	json, err := json.Marshal(res.Data)
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}

// Update modifies an existing user
func Update(us pb.UserServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing user param")
	}

	jsonStr := args[0]
	user := users.User{}

	err := json.Unmarshal([]byte(jsonStr), &user)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := us.Update(context.Background(), &pb.UpdateUserRequest{
		Data: user.ToProto(),
	})

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", fmt.Errorf("user with id = %v not found", user.ID)
		}

		return "", errors.New(err.Error())
	}

	json, err := json.Marshal(res.Data)
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}

// Delete removes a user by id
func Delete(us pb.UserServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing id param")
	}

	jsonStr := args[0]
	id := struct {
		ID string `json:"id"`
	}{}

	err := json.Unmarshal([]byte(jsonStr), &id)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := us.Delete(context.Background(), &pb.DeleteUserRequest{
		UserId: id.ID,
	})

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", fmt.Errorf("user with id = %v not found", id.ID)
		}

		return "", errors.New(err.Error())
	}

	json, err := json.Marshal(res.Data)
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}
