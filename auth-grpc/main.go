package main

import (
	"log"
	"net"

	"github.com/Edbeer/auth-grpc/pkg/db/psql"
	authpb "github.com/Edbeer/auth-grpc/proto"
	"github.com/Edbeer/auth-grpc/service"
	"github.com/Edbeer/auth-grpc/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	storage := storage.NewPostgresStorage(db)
	srv := service.NewAuthService(storage)
	
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))

	authpb.RegisterAuthServiceServer(server, srv)

	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("auth server start")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	server.GracefulStop()
}