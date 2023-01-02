package service

import (
	"context"

	authpb "github.com/Edbeer/auth-grpc/proto"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreateAccount(ctx context.Context, req *authpb.Account) (*authpb.Account, error) {
	return nil, nil
}

func (s *AuthService) GetAccount(req *authpb.GetAccountRequest, stream authpb.AuthService_GetAccountServer) error {
	return nil
}

func (s *AuthService) UpdateAccount(ctx context.Context, req *authpb.Account) (*authpb.Account, error) {
	return nil, nil
}

func (s *AuthService) DeleteAccount(ctx context.Context, req *authpb.Account) (*authpb.DeleteResponse, error) {
	return nil, nil
}

func (s *AuthService) GetAccountByID(ctx context.Context, req *authpb.Account) (*authpb.Account, error) {
	return nil, nil
}

func (s *AuthService) DepositAccount(ctx context.Context, req *authpb.Account) (*authpb.Account, error) {
	return nil, nil
}