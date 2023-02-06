//go:generate mockgen -source service.go -destination mock/storage_mock.go -package mock
package service

import (
	"context"
	"io"

	"github.com/Edbeer/auth-grpc/pkg/utils"
	authpb "github.com/Edbeer/payment-proto/auth-grpc/proto"
	"github.com/Edbeer/auth-grpc/types"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Storage interface {
	CreateAccount(ctx context.Context, account *authpb.CreateRequest) (*types.Account, error)
	UpdateAccount(ctx context.Context, account *authpb.UpdateRequest) (*types.Account, error)
	DeleteAccount(ctx context.Context, req *authpb.DeleteRequest) (*authpb.DeleteResponse, error)
	DepositAccount(ctx context.Context, req *authpb.DepositRequest) (*authpb.DepositResponse, error)
	GetAccountByID(ctx context.Context, req *authpb.GetIDRequest) (*types.Account, error)
	GetAccount(ctx context.Context) ([]*types.Account, error)
	SaveBalance(ctx context.Context, req *authpb.UpdateBalanceRequest) (*types.Account, error)
	UpdateStatement(ctx context.Context, req *authpb.StatementRequest) ([]string, error)
}

type RedisStorage interface {
	CreateSession(ctx context.Context, session *types.Session, expire int) (string, error)
	GetUserID(ctx context.Context, refreshToken string) (uuid.UUID, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	redisStorage RedisStorage
	storage      Storage
}

func NewAuthService(storage Storage, redisStorage RedisStorage) *AuthService {
	return &AuthService{
		storage:      storage,
		redisStorage: redisStorage,
	}
}

func (s *AuthService) CreateAccount(ctx context.Context, req *authpb.CreateRequest) (*authpb.AccountWithTokens, error) {
	accReq := &authpb.CreateRequest{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		CardNumber:       req.CardNumber,
		CardExpiryMonth:  req.CardExpiryMonth,
		CardExpiryYear:   req.CardExpiryYear,
		CardSecurityCode: req.CardSecurityCode,
	}

	account, err := s.storage.CreateAccount(ctx, accReq)
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.CreateJWT(account)
	if err != nil {
		return nil, err
	}

	// refreshToken
	refreshToken, err := s.redisStorage.CreateSession(ctx, &types.Session{
		UserID: account.ID,
	}, 86400)
	if err != nil {
		return nil, err
	}

	return accAndTokenToProto(account, accessToken, refreshToken), nil
}

func (s *AuthService) GetAccount(req *authpb.GetRequest, stream authpb.AuthService_GetAccountServer) error {
	accounts, err := s.storage.GetAccount(stream.Context())
	if err != nil {
		return err
	}
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			for _, account := range accounts {
				if err := stream.Send(accountToProto(account)); err != nil {
					return status.Error(codes.Canceled, "Stream has ended")
				}
			}
		}
		break
	}
	return nil
}

func (s *AuthService) UpdateAccount(ctx context.Context, req *authpb.UpdateRequest) (*authpb.Account, error) {
	updatedAccount, err := s.storage.UpdateAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	return accountToProto(updatedAccount), nil
}

func (s *AuthService) DeleteAccount(ctx context.Context, req *authpb.DeleteRequest) (*authpb.DeleteResponse, error) {
	status, err := s.storage.DeleteAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (s *AuthService) GetAccountByID(ctx context.Context, req *authpb.GetIDRequest) (*authpb.Account, error) {
	account, err := s.storage.GetAccountByID(ctx, req)
	if err != nil {
		return nil, err
	}
	return accountToProto(account), nil
}

func (s *AuthService) DepositAccount(ctx context.Context, req *authpb.DepositRequest) (*authpb.DepositResponse, error) {
	status, err := s.storage.DepositAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (s *AuthService) UpdateBalance(ctx context.Context, req *authpb.UpdateBalanceRequest) (*authpb.Account, error) {
	account, err := s.storage.SaveBalance(ctx, req)
	if err != nil {
		return nil, err
	}
	return accountToProto(account), nil
}

func (s *AuthService) CreateStatement(stream authpb.AuthService_CreateStatementServer) error {
	for i := 0; i < 2; i++ {
		statement, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = s.storage.UpdateStatement(stream.Context(), statement)
		if err != nil {
			return err
		}
		if err = stream.Send(&authpb.StatementResponse{}); err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) GetStatement(req *authpb.StatementGet, stream authpb.AuthService_GetStatementServer) error {
	account, err := s.storage.GetAccountByID(stream.Context(), &authpb.GetIDRequest{
		Id: req.AccountId,
	})
	if err != nil {
		return err
	}
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			for _, statement := range account.Statement {
				if err := stream.Send(&authpb.Statement{
					PaymentId: statement,
				}); err != nil {
					return err
				}
			}
		}
		break
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, req *authpb.LoginRequest) (*authpb.AccountWithTokens, error) {
	account, err := s.storage.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	accessToken, err := utils.CreateJWT(account)
	if err != nil {
		return nil, err
	}

	// refreshToken
	refreshToken, err := s.redisStorage.CreateSession(ctx, &types.Session{
		UserID: account.ID,
	}, 86400)
	if err != nil {
		return nil, err
	}

	return accAndTokenToProto(account, accessToken, refreshToken), nil
}

func (s *AuthService) SignOut(ctx context.Context, req *authpb.QuitRequest) (*authpb.QuitResponse, error) {
	err := s.redisStorage.DeleteSession(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}	
	return &authpb.QuitResponse{
		Message: "sign-out",
	}, err	
}

func (s *AuthService) RefreshTokens(ctx context.Context, req *authpb.RefreshRequest) (*authpb.Tokens, error) {
	uid, err := s.redisStorage.GetUserID(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	account, err := s.storage.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: uid.String(),
	})
	if err != nil {
		return nil, err
	}

	// jwt-token
	tokenString, err := utils.CreateJWT(account)
	if err != nil {
		return nil, err
	}
	// refreshToken
	refreshToken, err := s.redisStorage.CreateSession(ctx, &types.Session{
		UserID: account.ID,
	}, 86400)
	if err != nil {
		return nil, err
	}

	return &authpb.Tokens{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}, nil
}

func accountToProto(acc *types.Account) *authpb.Account {
	return &authpb.Account{
		Id:               acc.ID.String(),
		FirstName:        acc.FirstName,
		LastName:         acc.LastName,
		CardNumber:       acc.CardNumber,
		CardExpiryMonth:  acc.CardExpiryMonth,
		CardExpiryYear:   acc.CardExpiryYear,
		CardSecurityCode: acc.CardSecurityCode,
		Balance:          acc.Balance,
		BlockedMoney:     acc.BlockedMoney,
		CreatedAt:        timestamppb.New(acc.CreatedAt),
	}
}

func accAndTokenToProto(acc *types.Account, accessToken, refreshToken string) *authpb.AccountWithTokens {
	return &authpb.AccountWithTokens{
		Account: &authpb.Account{
			Id:               acc.ID.String(),
			FirstName:        acc.FirstName,
			LastName:         acc.LastName,
			CardNumber:       acc.CardNumber,
			CardExpiryMonth:  acc.CardExpiryMonth,
			CardExpiryYear:   acc.CardExpiryYear,
			CardSecurityCode: acc.CardSecurityCode,
			Balance:          acc.Balance,
			BlockedMoney:     acc.BlockedMoney,
			CreatedAt:        timestamppb.New(acc.CreatedAt),
		},
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}
}
