package main

import (
	"log"
	"net"

	authpb "github.com/Edbeer/auth-grpc/proto"
	"github.com/Edbeer/auth-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	server := grpc.NewServer()
	srv := service.NewAuthService()

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