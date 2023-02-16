package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Edbeer/api-gateway/pkg/utils"
	authpb "github.com/Edbeer/payment-proto/auth-grpc/proto"
)

type CreateRequest struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	CardNumber       string `json:"card_number"`
	CardExpiryMonth  string `json:"card_expiry_month"`
	CardExpiryYear   string `json:"card_expiry_year"`
	CardSecurityCode string `json:"card_security_code"`
}

// createAccount godoc
// @Summary Create new account
// @Description register new account, returns account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body CreateRequest true "create account info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account [post]
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

// refreshTokens godoc
// @Summary Refresh tokens
// @Description refresh access and refresh tokens, returns tokens
// @Tags Account
// @Accept json
// @Produce json
// @Param input body RefreshRequest true "refresh tokens account info"
// @Success 200 {object} Tokens
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/refresh [post]
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

// signIn godoc
// @Summary Login
// @Description log in to your account, returns account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body LoginRequest true "login account info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/sign-in [post]
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


// signOut godoc
// @Summary Logout
// @Description log out of your account, returns status
// @Tags Account
// @Produce json
// @Success 200 {integer} 200
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/sign-out [post]
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

// updateAccount godoc
// @Summary Update account
// @Description update account, returns updated account
// @Tags Account
// @Accept json
// @Produce json
// @Param id path string true "update account info"
// @Param input body UpdateRequest true "update account info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/{id} [put]
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

// deleteAccount godoc
// @Summary Delete account
// @Description delete account, returns status
// @Tags Account
// @Produce json
// @Param id path string true "delete account info"
// @Success 200 {integer} 200
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/{id} [delete]
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

// depositAccount godoc
// @Summary Deposit money
// @Description deposit money to account, returns account
// @Tags Account
// @Accept json
// @Produce json
// @Param input body DepositRequest true "deposit account info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/deposit [post]
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

// getAccountByID godoc
// @Summary Get account by id
// @Description get account by id, returns account
// @Tags Account
// @Produce json
// @Param id path string true "get account by id info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/{id} [get]
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

// getAccount godoc
// @Summary Get all accounts
// @Description get all accounts, returns accounts
// @Tags Account
// @Produce json
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account [get]
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

// getStatement godoc
// @Summary Get account statement
// @Description get account statement, returns statement
// @Tags Account
// @Produce json
// @Param id path string true "get statement info"
// @Failure 400  {object}  utils.ApiError
// @Failure 404  {object}  utils.ApiError
// @Failure 500  {object}  utils.ApiError
// @Router /account/statement/{id} [get]
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
