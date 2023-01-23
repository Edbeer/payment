package auth

import (
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/auth/routes"
	"github.com/Edbeer/api-gateway/pkg/auth/utils"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router) *AuthClient {
	client := &AuthClient{
		client: AuthServiceClient(),
	}
	// POST
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/account", HTTPHandler(client.CreateAccount))
	postRouter.HandleFunc("/account/deposit", HTTPHandler(client.DepositAccount))
	// GET
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/account", HTTPHandler(client.GetAccount))
	getRouter.HandleFunc("/account/{id}", HTTPHandler(client.GetAccountByID))
	getRouter.HandleFunc("/account/statement/{id}", HTTPHandler(client.GetStatement))
	// PUT
	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/account/{id}", HTTPHandler(client.UpdateAccount))
	// DELETE
	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/account/{id}", HTTPHandler(client.DeleteAccount))
	return client
}

// Create Account
func (s *AuthClient) CreateAccount(w http.ResponseWriter, r *http.Request) error {
	return routes.CreateAccount(w, r, s.client)
}

// Deposit Account
func (s *AuthClient) DepositAccount(w http.ResponseWriter, r *http.Request) error {
	return routes.DepositAccount(w, r, s.client)
}

// Get All Accounts
func (s *AuthClient) GetAccount(w http.ResponseWriter, r *http.Request) error {
	return routes.GetAccount(w, r, s.client)
}

// Get Account By ID
func (s *AuthClient) GetAccountByID(w http.ResponseWriter, r *http.Request) error {
	return routes.GetAccountByID(w, r, s.client)
}

// Get all Statement
func (s *AuthClient) GetStatement(w http.ResponseWriter, r *http.Request) error {
	return routes.GetStatement(w, r, s.client)
}

// Update Account
func (s *AuthClient) UpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return routes.UpdateAccount(w, r, s.client)
}

// Delete Account
func (s *AuthClient) DeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return routes.DeleteAccount(w, r, s.client)
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

// Wrapper for handler func
func HTTPHandler(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
		}
	}
}
