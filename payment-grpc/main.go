package main

import (
	"log"
	"net"

	authpb "github.com/Edbeer/payment-proto/auth-grpc/proto"
	"github.com/Edbeer/payment-grpc/pkg/db"
	paymentpb "github.com/Edbeer/payment-proto/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/service"
	"github.com/Edbeer/payment-grpc/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// postgres db
	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("init postgres")

	// postgres storage
	storage := storage.NewPostgresStorage(db)

	// client
	conn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := authpb.NewAuthServiceClient(conn)

	// payment service
	srv := service.NewPaymentService(storage, client, db)
	// grpc server
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))
	// register service
	paymentpb.RegisterPaymentServiceServer(server, srv)
	// reflection 
	reflection.Register(server)
	// listen on port :50051
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
