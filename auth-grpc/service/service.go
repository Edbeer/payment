//go:generate mockgen -source service.go -destination mock/storage_mock.go -package mock
package service

import (
	"context"
	"io"

	authpb "github.com/Edbeer/auth-grpc/proto"
	"github.com/Edbeer/auth-grpc/types"
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

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
	storage Storage
}

func NewAuthService(storage Storage) *AuthService {
	return &AuthService{
		storage: storage,
	}
}

func (s *AuthService) CreateAccount(ctx context.Context, req *authpb.CreateRequest) (*authpb.Account, error) {
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

	return accountToProto(account), nil
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
				if err := stream.SendMsg(accountToProto(account)); err != nil {
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
	for {
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

// func (s *AuthService) CreateStatement(ctx context.Context, req *authpb.StatementRequest) (*authpb.StatementResponse, error) {
// 	_, err := s.storage.UpdateStatement(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &authpb.StatementResponse{}, nil
// }

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
