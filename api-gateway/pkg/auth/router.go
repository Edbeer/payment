package auth

import (
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/auth/routes"
	"github.com/Edbeer/api-gateway/pkg/utils"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router) *AuthClient {
	client := &AuthClient{
		client: AuthServiceClient(),
	}
	// POST
	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/account", utils.HTTPHandler(client.CreateAccount))
	postRouter.HandleFunc("/account/deposit", utils.HTTPHandler(client.DepositAccount))
	postRouter.HandleFunc("/account/sign-in", utils.HTTPHandler(client.SignIn))
	postRouter.HandleFunc("/account/sign-out", utils.HTTPHandler(client.Signout))
	postRouter.HandleFunc("/account/refresh", utils.HTTPHandler(client.RefreshTokens))
	// GET
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/account", utils.HTTPHandler(client.GetAccount))
	getRouter.HandleFunc("/account/{id}", AuthJWT(utils.HTTPHandler(client.GetAccountByID)))
	getRouter.HandleFunc("/account/statement/{id}", AuthJWT(utils.HTTPHandler(client.GetStatement)))
	// PUT
	putRouter := router.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/account/{id}", AuthJWT(utils.HTTPHandler(client.UpdateAccount)))
	// DELETE
	deleteRouter := router.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/account/{id}", AuthJWT(utils.HTTPHandler(client.DeleteAccount)))
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

// SignIn
func (s *AuthClient) SignIn(w http.ResponseWriter, r *http.Request) error {
	return routes.SignIn(w, r, s.client)
}

// SignOut
func (s *AuthClient) Signout(w http.ResponseWriter, r *http.Request) error {
	return routes.SignOut(w, r, s.client)
}

// RefreshTokens
func (s *AuthClient) RefreshTokens(w http.ResponseWriter, r *http.Request) error {
	return routes.RefreshTokens(w, r, s.client)
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
