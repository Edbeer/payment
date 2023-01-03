package types

import (
	authpb "github.com/Edbeer/auth-grpc/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateAuthPayment(req *paymentpb.Payment, customer *authpb.Account, merchant *authpb.Account) *paymentpb.Payment {
	mid, err := uuid.Parse(merchant.Id)
	if err != nil {
		return nil
	}
	cid, err := uuid.Parse(customer.Id)
	if err != nil {
		return nil
	}
	return &paymentpb.Payment{
		PaymentId:       uuid.NewString(),
		PaymentReceiver: mid.String(),
		Payer:           cid.String(),
		Currency:        req.Currency,
		Operation:       req.Operation,
		Status:          req.Status,
		Amount:          req.Amount,
		CreatedAt:       timestamppb.Now(),
	}
}
