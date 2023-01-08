package types

import (
	"time"

	authpb "github.com/Edbeer/auth-grpc/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/google/uuid"
)

type Payment struct {
	PaymentId       uuid.UUID `json:"id"`
	Merchant        uuid.UUID `json:"merchant"`
	Customer        uuid.UUID `json:"customer"`
	CardNumber      string    `json:"card_number"`
	CardExpiryMonth string    `json:"card_expiry_month"`
	CardExpiryYear  string    `json:"card_expiry_year"`
	Currency        string    `json:"currency"`
	Operation       string    `json:"operation"`
	Status          string    `json:"status"`
	Amount          uint64    `json:"amount"`
	CreatedAt       time.Time `json:"creation_at"`
}

func CreateAuthPayment(req *paymentpb.CreateRequest, customer *authpb.Account, merchant *authpb.Account, status string) *Payment {
	mid, err := uuid.Parse(merchant.Id)
	if err != nil {
		return nil
	}
	cid, err := uuid.Parse(customer.Id)
	if err != nil {
		return nil
	}
	return &Payment{
		PaymentId:       uuid.New(),
		Merchant:        mid,
		Customer:        cid,
		CardNumber:      req.CardNumber,
		CardExpiryMonth: req.CardExpiryMonth,
		CardExpiryYear:  req.CardExpiryYear,
		Currency:        req.Currency,
		Operation:       req.Operation,
		Status:          status,
		Amount:          req.Amount,
		CreatedAt:       time.Now(),
	}
}

type Statement struct {
	PaymentId   string `json:"payment_id"`
	AccountId  string `json:"account_id"`
}
