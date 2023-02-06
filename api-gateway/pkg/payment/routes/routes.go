package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/utils"
	paymentpb "github.com/Edbeer/payment-proto/payment-grpc/proto"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Merchant        uuid.UUID `json:"merchant_id"`
	Customer        uuid.UUID `json:"customer_id"`
	CardNumber       string    `json:"card_number"`
	CardExpiryMonth  string    `json:"card_expiry_month"`
	CardExpiryYear   string    `json:"card_expiry_year"`
	CardSecurityCode string    `json:"card_security_code"`
	Currency         string    `json:"currency"`
	Amount           uint64    `json:"amount"`
}

func CreatePayment(w http.ResponseWriter, r *http.Request, cc paymentpb.PaymentServiceClient) error {
	req := &CreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	statement, err := cc.CreatePayment(r.Context(), &paymentpb.CreateRequest{
		Merchant:         req.Merchant.String(),
		Customer:         req.Customer.String(),
		CardNumber:       req.CardNumber,
		CardExpiryMonth:  req.CardExpiryMonth,
		CardExpiryYear:   req.CardExpiryYear,
		CardSecurityCode: req.CardSecurityCode,
		Currency:         req.Currency,
		Amount:           req.Amount,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, statement)
}

type PaidRequest struct {
	Amount    uint64    `json:"amount"`
}

func CapturePayment(w http.ResponseWriter, r *http.Request, cc paymentpb.PaymentServiceClient) error {
	uuid, err := utils.GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	
	req := &PaidRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	statement, err := cc.CapturePayment(r.Context(), &paymentpb.PaidRequest{
		PaymentId: uuid.String(),
		Amount:    req.Amount,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, statement)
}

func CancelPayment(w http.ResponseWriter, r *http.Request, cc paymentpb.PaymentServiceClient) error {
	uuid, err := utils.GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	
	req := &PaidRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	statement, err := cc.CancelPayment(r.Context(), &paymentpb.PaidRequest{
		PaymentId: uuid.String(),
		Amount:    req.Amount,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, statement)
}

func RefundPayment(w http.ResponseWriter, r *http.Request, cc paymentpb.PaymentServiceClient) error {
	uuid, err := utils.GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	
	req := &PaidRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	statement, err := cc.RefundPayment(r.Context(), &paymentpb.PaidRequest{
		PaymentId: uuid.String(),
		Amount:    req.Amount,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, statement)
}