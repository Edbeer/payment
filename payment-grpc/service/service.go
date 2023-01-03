package service

import (
	"context"

	paymentpb "github.com/Edbeer/payment-grpc/proto"
)

type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}

func (s *PaymentService) CapturePayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}

func (s *PaymentService) CancelPayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}

func (s *PaymentService) GetStatemet(req *paymentpb.StatementRequest, stream paymentpb.PaymentService_GetStatemetServer) error {
	return nil
}