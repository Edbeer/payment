package main

import (
	"log"
	"net"

	"github.com/Edbeer/payment-grpc/storage"
	"github.com/Edbeer/payment-grpc/pkg/db"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewPostgresStorage(db)

	srv := service.NewPaymentService(storage)
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))
	
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
