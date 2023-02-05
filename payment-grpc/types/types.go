package types

import (
	"time"

	authpb "github.com/Edbeer/auth/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/google/uuid"
)

type Payment struct {
	PaymentId       uuid.UUID `json:"payment_id"`
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
		Operation:       "Authorization",
		Status:          status,
		Amount:          req.Amount,
		CreatedAt:       time.Now(),
	}
}

// creating a complete payment
func CreateCompletePayment(paidReq *paymentpb.PaidRequest, referncedPayment *Payment, status string) *Payment {
	return &Payment{
		PaymentId:       uuid.New(),
		Merchant:        referncedPayment.Merchant,
		Customer:        referncedPayment.Customer,
		CardNumber:      referncedPayment.CardNumber,
		CardExpiryMonth: referncedPayment.CardExpiryMonth,
		CardExpiryYear:  referncedPayment.CardExpiryYear,
		Currency:        referncedPayment.Currency,
		Operation:       referncedPayment.Operation,
		Status:          status,
		Amount:          paidReq.Amount,
		CreatedAt:       time.Now(),
	}
}

