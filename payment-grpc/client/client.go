package main

import (
	"context"
	"io"
	"log"

	authpb "github.com/Edbeer/auth-grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := authpb.NewAuthServiceClient(conn)

	// createAccount(client)
	// getAccount(client)
	getAccountByID(client)
}

func createAccount(client authpb.AuthServiceClient) {
	r1 := &authpb.CreateRequest{
		FirstName:        "Pasha",
		LastName:         "Volkov",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "12",
		CardSecurityCode: "123",
	}
	r2 := &authpb.CreateRequest{
		FirstName:        "Pasha1",
		LastName:         "Volkov1",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "12",
		CardSecurityCode: "123",
	}

	resp1, err := client.CreateAccount(context.Background(), r1)
	if err != nil {
		log.Fatal(err)
	}
	resp2, _ := client.CreateAccount(context.Background(), r2)
	log.Println(resp1, resp2)
}

func getAccount(client authpb.AuthServiceClient) {
	req := &authpb.GetRequest{}
	stream, err := client.GetAccount(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	for {
		value, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(value)
	}
}

func getAccountByID(client authpb.AuthServiceClient) {
	stream, err := client.GetAccountByID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	accs := []*authpb.GetIDRequest{
		{
			Id: "48ed1fd7-db48-49d1-9869-d01349a8d4a7",
		},
		{
			Id: "f36f7969-817b-430f-9a60-ce0c01e4d4eb",
		},
	}
	for _, acc := range accs {
		if err := stream.Send(acc); err != nil {
			log.Fatal(err)
		}
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatal(err)
	}
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(value)
	}
}
