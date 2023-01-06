package service

import (
	"context"
	"database/sql"

	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/types"
)

type Storage interface {
	SavePayment(ctx context.Context, tx *sql.Tx, payment *types.Payment) (*types.Payment, error)
}

type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
	storage Storage
}

func NewPaymentService(storage Storage) *PaymentService {
	return &PaymentService{storage: storage}
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