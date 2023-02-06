package payment

import (
	paymentpb "github.com/Edbeer/payment-proto/payment-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentClient struct {
	client paymentpb.PaymentServiceClient
}

func PaymentServiceClient() paymentpb.PaymentServiceClient {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}

	return paymentpb.NewPaymentServiceClient(conn)
}