package service

import (
	"context"
	"log"

	authpb "github.com/Edbeer/auth-grpc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthService struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) CreateAccount(ctx context.Context, req *authpb.Account) (*authpb.Account, error) {
	log.Println("Geeee")
	return &authpb.Account{
		Id:               "",
		FirstName:        req.FirstName,
		LastName:         "",
		CardNumber:       "",
		CardExpiryMonth:  "",
		CardExpiryYear:   "",
		CardSecurityCode: "",
		Balance:          0,
		BlockedMoney:     0,
		Statement:        []*authpb.Statement{},
		CreatedAt:        timestamppb.Now(),
	}, nil
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

func (s *AuthService) GetAccountByID(stream authpb.AuthService_GetAccountByIDServer) error {
	return nil
}

func (s *AuthService) DepositAccount(ctx context.Context, req *authpb.Account) (*authpb.Account, error) {
	return nil, nil
}