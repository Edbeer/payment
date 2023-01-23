package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/auth/utils"
	authpb "github.com/Edbeer/auth-grpc/proto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	CardNumber       string `json:"card_number"`
	CardExpiryMonth  string `json:"card_expiry_month"`
	CardExpiryYear   string `json:"card_expiry_year"`
	CardSecurityCode string `json:"card_security_code"`
}

func CreateAccount(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	req := &CreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	account, err := cc.CreateAccount(r.Context(), &authpb.CreateRequest{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		CardNumber:       req.CardNumber,
		CardExpiryMonth:  req.CardExpiryMonth,
		CardExpiryYear:   req.CardExpiryYear,
		CardSecurityCode: req.CardSecurityCode,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	// TODO Cookie

	return utils.WriteJSON(w, http.StatusOK, account)
}

type UpdateRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	CardNumber       string `json:"card_number"`
	CardExpiryMonth  string `json:"card_expiry_month"`
	CardExpiryYear   string `json:"card_expiry_year"`
	CardSecurityCode string `json:"card_security_code"`
}

func UpdateAccount(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	uuid, err := GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	
	req := &UpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	account, err := cc.UpdateAccount(r.Context(), &authpb.UpdateRequest{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		CardNumber:       req.CardNumber,
		CardExpiryMonth:  req.CardExpiryMonth,
		CardExpiryYear:   req.CardExpiryYear,
		CardSecurityCode: req.CardSecurityCode,
		Id:               uuid.String(),
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, account)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	uuid, err := GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	_, err = cc.DeleteAccount(r.Context(), &authpb.DeleteRequest{
		Id: uuid.String(),
	})

	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, "account was deleted")
}

type DepositRequest struct {
	CardNumber string `json:"card_number"`
	Balance    uint64 `json:"balance"`
}

func DepositAccount(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	req := &DepositRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	status, err := cc.DepositAccount(r.Context(), &authpb.DepositRequest{
		CardNumber: req.CardNumber,
		Balance:    req.Balance,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, status)
}

func GetAccountByID(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	uuid, err := GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	account, err := cc.GetAccountByID(r.Context(), &authpb.GetIDRequest{
		Id: uuid.String(),
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, account)
}

type UpdateBalanceRequest struct {
	Id           uuid.UUID `json:"id"`
	Balance      uint64    `json:"balance"`
	BlockedMoney uint64    `json:"blocked_money"`
}

// TODO remove
func UpdateBalance(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	req := &UpdateBalanceRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	account, err := cc.UpdateBalance(r.Context(), &authpb.UpdateBalanceRequest{
		Id:           req.Id.String(),
		Balance:      req.Balance,
		BlockedMoney: req.BlockedMoney,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	return utils.WriteJSON(w, http.StatusOK, account)
}

func GetAccount(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	stream, err := cc.GetAccount(r.Context(), &authpb.GetRequest{})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	accounts := []*authpb.Account{}
	for {
		account, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
		}

		accounts = append(accounts, account)
	}

	return utils.WriteJSON(w, http.StatusOK, accounts)
}

func GetStatement(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	uuid, err := GetUUID(r)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	stream, err := cc.GetStatement(r.Context(), &authpb.StatementGet{
		AccountId: uuid.String(),
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	statements := []*authpb.Statement{}
	for {
		statement, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
		}
		statements = append(statements, statement)
	}

	return utils.WriteJSON(w, http.StatusOK, statements)
}

// Get id from url
func GetUUID(r *http.Request) (uuid.UUID, error) {
	id := mux.Vars(r)["id"]
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}
