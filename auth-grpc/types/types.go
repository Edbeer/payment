package types

import (
	"time"

	authpb "github.com/Edbeer/auth-grpc/proto"
	"github.com/google/uuid"
)

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