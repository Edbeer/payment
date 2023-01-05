package main

import (

	"log"
	"net"

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


	log.Println("payment server start")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	server.GracefulStop()
}
