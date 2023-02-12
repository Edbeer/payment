package main

import (
	"log"
	"net"

	"github.com/Edbeer/auth-grpc/pkg/db/psql"
	red "github.com/Edbeer/auth-grpc/pkg/db/redis"
	authpb "github.com/Edbeer/payment-proto/auth-grpc/proto"
	"github.com/Edbeer/auth-grpc/service"
	"github.com/Edbeer/auth-grpc/storage/psql"
	"github.com/Edbeer/auth-grpc/storage/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// postgres db
	db, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("init postgres")
	// redis
	redisClient := red.NewRedisClient()
	defer redisClient.Close()
	log.Println("init redis")

	storage := storage.NewPostgresStorage(db)
	redis := redisrepo.NewRedisStorage(redisClient)

	srv := service.NewAuthService(storage, redis)
	
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