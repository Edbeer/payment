package types

import (
	"time"

	authpb "github.com/Edbeer/auth-grpc/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/google/uuid"
)

type Payment struct {
	PaymentId       uuid.UUID            `json:"id"`
	PaymentReceiver uuid.UUID            `json:"payment_receiver"`
	Payer           uuid.UUID            `json:"payer"`
	Currency        paymentpb.Currencies `json:"currency"`
	Operation       paymentpb.Operations `json:"operation"`
	Status          paymentpb.Statuses   `json:"status"`
	Amount          float64              `json:"amount"`
	CreatedAt       time.Time            `json:"creation_at"`
}

func CreateAuthPayment(req *paymentpb.Payment, customer *authpb.Account, merchant *authpb.Account) *Payment {
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
		PaymentReceiver: mid,
		Payer:           cid,
		Currency:        req.Currency,
		Operation:       req.Operation,
		Status:          req.Status,
		Amount:          req.Amount,
		CreatedAt:       time.Now(),
	}
}

// Account
type Account struct {
	ID               uuid.UUID `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	CardNumber       string    `json:"card_number"`
	CardExpiryMonth  string    `json:"card_expiry_month"`
	CardExpiryYear   string    `json:"card_expiry_year"`
	CardSecurityCode string    `json:"card_security_code"`
	Balance          uint64    `json:"balance"`
	BlockedMoney     uint64    `json:"blocked_money"`
	CreatedAt        time.Time `json:"created_at"`
}

func NewAccount(req *authpb.CreateRequest) *Account {
	return &Account{
		ID:               uuid.New(),
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		CardNumber:       req.CardNumber,
		CardExpiryMonth:  req.CardExpiryMonth,
		CardExpiryYear:   req.CardExpiryYear,
		CardSecurityCode: req.CardSecurityCode,
		Balance:          0,
		BlockedMoney:     0,
		CreatedAt:        time.Now(),
	}
}
