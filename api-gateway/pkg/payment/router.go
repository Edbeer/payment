package payment

import (
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/payment/routes"
	"github.com/Edbeer/api-gateway/pkg/utils"
	"github.com/gorilla/mux"
)

func RegisterPaymentRouter(router *mux.Router) *PaymentClient {
	client := &PaymentClient{
		client: PaymentServiceClient(),
	}

	// POST
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/payment/auth", utils.HTTPHandler(client.CreatePayment))
	postRouter.HandleFunc("/payment/capture/{id}", utils.HTTPHandler(client.CapturePayment))
	postRouter.HandleFunc("/payment/cancel/{id}", utils.HTTPHandler(client.CancelPayment))
	postRouter.HandleFunc("/payment/refund/{id}", utils.HTTPHandler(client.RefundPayment))

	return client
}

func (s *PaymentClient) CreatePayment(w http.ResponseWriter, r *http.Request) error {
	return routes.CreatePayment(w, r, s.client)
}

func (s *PaymentClient) CapturePayment(w http.ResponseWriter, r *http.Request) error {
	return routes.CapturePayment(w, r, s.client)
}

func (s *PaymentClient) CancelPayment(w http.ResponseWriter, r *http.Request) error {
	return routes.CapturePayment(w, r, s.client)
}

func (s *PaymentClient) RefundPayment(w http.ResponseWriter, r *http.Request) error {
	return routes.CapturePayment(w, r, s.client)
}