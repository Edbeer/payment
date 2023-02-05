package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/utils"
	authpb "github.com/Edbeer/auth/proto"
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

	accountWithToken, err := cc.CreateAccount(r.Context(), &authpb.CreateRequest{
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
	w.Header().Add("x-jwt-token", accountWithToken.AccessToken)
	// cookie
	cookie := &http.Cookie{
		Name:       "refresh-token",
		Value:      accountWithToken.RefreshToken,
		Path:       "/",
		RawExpires: "",
		MaxAge:     86400,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
	}
	http.SetCookie(w, cookie)

	return utils.WriteJSON(w, http.StatusOK, accountWithToken.Account)
}

type RefreshRequest struct {
	RefreshRequest string `json:"refresh_token"`
}

type Tokens struct {
	RefreshRequest string `json:"refresh_token"`
	AccessToken    string `json:"access_token"`
}

func RefreshTokens(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	req := &RefreshRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	tokens, err := cc.RefreshTokens(r.Context(), &authpb.RefreshRequest{
		RefreshToken: req.RefreshRequest,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	w.Header().Add("x-jwt-token", tokens.AccessToken)
	// cookie
	cookie := &http.Cookie{
		Name:       "refresh-token",
		Value:      tokens.RefreshToken,
		Path:       "/",
		RawExpires: "",
		MaxAge:     86400,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
	}
	http.SetCookie(w, cookie)

	return utils.WriteJSON(w, http.StatusOK, &Tokens{
		RefreshRequest: tokens.RefreshToken,
		AccessToken:    tokens.AccessToken,
	})
}

type LoginRequest struct {
	Id string `json:"id"`
}

func SignIn(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	req := &LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}

	accountWithToken, err := cc.SignIn(r.Context(), &authpb.LoginRequest{
		Id: req.Id,
	})
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	w.Header().Add("x-jwt-token", accountWithToken.AccessToken)
	// cookie
	cookie := &http.Cookie{
		Name:       "refresh-token",
		Value:      accountWithToken.RefreshToken,
		Path:       "/",
		RawExpires: "",
		MaxAge:     86400,
		Secure:     false,
		HttpOnly:   true,
		SameSite:   0,
	}
	http.SetCookie(w, cookie)

	return utils.WriteJSON(w, http.StatusOK, accountWithToken.Account)
}

func SignOut(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		if err == http.ErrNoCookie {
			return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: "Cookie doesn't exist"})
		}
		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
	}
	message, err := cc.SignOut(r.Context(), &authpb.QuitRequest{
		RefreshToken: cookie.Value,
	})

	return utils.WriteJSON(w, http.StatusOK, message.Message)
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
	uuid, err := utils.GetUUID(r)
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
	uuid, err := utils.GetUUID(r)
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
	uuid, err := utils.GetUUID(r)
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

// type UpdateBalanceRequest struct {
// 	Id           uuid.UUID `json:"id"`
// 	Balance      uint64    `json:"balance"`
// 	BlockedMoney uint64    `json:"blocked_money"`
// }

// // TODO remove
// func UpdateBalance(w http.ResponseWriter, r *http.Request, cc authpb.AuthServiceClient) error {
// 	req := &UpdateBalanceRequest{}

// 	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
// 		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
// 	}

// 	account, err := cc.UpdateBalance(r.Context(), &authpb.UpdateBalanceRequest{
// 		Id:           req.Id.String(),
// 		Balance:      req.Balance,
// 		BlockedMoney: req.BlockedMoney,
// 	})
// 	if err != nil {
// 		return utils.WriteJSON(w, http.StatusBadRequest, utils.ApiError{Error: err.Error()})
// 	}

// 	return utils.WriteJSON(w, http.StatusOK, account)
// }

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
	uuid, err := utils.GetUUID(r)
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
