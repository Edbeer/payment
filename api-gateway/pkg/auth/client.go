package auth

import (
	authpb "github.com/Edbeer/payment-proto/auth-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client authpb.AuthServiceClient
}

func AuthServiceClient() authpb.AuthServiceClient {
	conn, err := grpc.Dial(":50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	return authpb.NewAuthServiceClient(conn)
}