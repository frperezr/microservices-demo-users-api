package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/frperezr/noken-test/pb"

	"github.com/frperezr/noken-test/src/users-api/database"
	userService "github.com/frperezr/noken-test/src/users-api/rpc/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	port := os.Getenv("PORT")
	postgresDSN := os.Getenv("POSTGRES_DSN")

	if port == "" {
		log.Fatal("missing env variable PORT")
	}

	if postgresDSN == "" {
		log.Fatal("missing env variable POSTGRES_DSN")
	}

	postgresService, err := database.NewPostgres(postgresDSN)
	if err != nil {
		log.Fatalf("Failed connect to postgres: %v", err)
	}

	server := grpc.NewServer()
	service := userService.New(postgresService)

	pb.RegisterUserServiceServer(server, service)
	reflection.Register(server)

	log.Println("Starting User service...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Failed to list: %v", err)
	}

	log.Println(fmt.Sprintf("User service, Listening on: %v", port))

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Fatal to serve: %v", err)
	}
}
