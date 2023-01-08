package service

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	authpb "github.com/Edbeer/auth-grpc/proto"
	mockstore "github.com/Edbeer/auth-grpc/service/mock"
	"github.com/Edbeer/auth-grpc/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage)

	req := &authpb.CreateRequest{
		FirstName:        "Pasha",
		LastName:         "Volkov",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
	}
	reqAcc := types.NewAccount(req)
	acc := &types.Account{
		ID:               reqAcc.ID,
		FirstName:        "Pasha",
		LastName:         "Volkov",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          0,
		BlockedMoney:     0,
		CreatedAt:        time.Now(),
	}
	mockStorage.EXPECT().CreateAccount(context.Background(), gomock.Eq(req)).Return(acc, nil).AnyTimes()

	account, err := mockService.CreateAccount(context.Background(), req)
	require.NoError(t, err)
	require.Nil(t, err)
	require.Equal(t, accountToProto(acc), account)
}

func Test_UpdateAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage)

	uid := uuid.New()
	reqToUpdate := &authpb.UpdateRequest{
		FirstName:        "Pasha1",
		LastName:         "",
		CardNumber:       "4444444444444443",
		CardExpiryMonth:  "",
		CardExpiryYear:   "",
		CardSecurityCode: "",
		Id:               uid.String(),
	}

	acc := &types.Account{
		ID:               uid,
		FirstName:        "Pasha1",
		LastName:         "Volkov",
		CardNumber:       "4444444444444443",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          0,
		BlockedMoney:     0,
		CreatedAt:        time.Now(),
	}

	mockStorage.EXPECT().UpdateAccount(context.Background(), reqToUpdate).Return(acc, nil).AnyTimes()
	require.Equal(t, reqToUpdate.FirstName, acc.FirstName)
	require.Equal(t, reqToUpdate.CardNumber, acc.CardNumber)
	
	account, err := mockService.UpdateAccount(context.Background(), reqToUpdate)
	require.NoError(t, err)
	require.Equal(t, account, accountToProto(acc))
}

func Test_DeleteAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage)

	req := &authpb.DeleteRequest{
		Id: uuid.New().String(),
	}
	resp := &authpb.DeleteResponse{
		Status: "Account was deleted",
	}
	mockStorage.EXPECT().DeleteAccount(context.Background(), req).Return(resp, nil).AnyTimes()

	result, err := mockService.DeleteAccount(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, result, resp)
}

func Test_DepositAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage)

	reqDep := &authpb.DepositRequest{
		CardNumber: "4444444444444444",
		Balance:    50,
	}

	resp := &authpb.DepositResponse{
		Status: "Successful deposit",
	}

	mockStorage.EXPECT().DepositAccount(context.Background(), reqDep).Return(resp, nil).AnyTimes()

	result, err := mockService.DepositAccount(context.Background(), reqDep)
	require.NoError(t, err)
	require.Equal(t, result, resp)
}
