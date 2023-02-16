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

// createPayment godoc
// @Summary Create payment
// @Description Create payment: Acceptance of payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path string true "create payment info"
// @Param input body CreateRequest true "create payment info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /payment/auth [post]
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

// capturePayment godoc
// @Summary Capture payment
// @Description Capture payment: Successful payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path string true "capture payment info"
// @Param input body PaidRequest true "capture payment info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /payment/capture/{id} [post]
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

// cancelPayment godoc
// @Summary Cancel payment
// @Description Cancel payment: cancel authorization payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path string true "cancel payment info"
// @Param input body PaidRequest true "cancel payment info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /payment/cancel/{id} [post]
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

// refundPayment godoc
// @Summary Refund payment
// @Description Refund: Refunded payment, if there is a refund
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path string true "refund payment info"
// @Param input body PaidRequest true "refund payment info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /payment/refund/{id} [post]
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