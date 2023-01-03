package main

import (
	"context"
	"log"
	"net"

	authpb "github.com/Edbeer/auth-grpc/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))
	srv := service.NewPaymentService()
	
	paymentpb.RegisterPaymentServiceServer(server, srv)

	reflection.Register(server)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	
	client := authpb.NewAuthServiceClient(conn)
	req := &authpb.Account{
		FirstName:        "Pasha",
	}
	res, err := client.CreateAccount(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
	
	log.Println("payment server start")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	server.GracefulStop()
}