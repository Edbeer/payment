package service

import (
	"context"
	"fmt"

	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mockclient "github.com/Edbeer/auth/client/mock"
	authpb "github.com/Edbeer/auth/proto"
	mockstore "github.com/Edbeer/auth/service/mock"
	"github.com/Edbeer/auth/types"
	// "github.com/alicebob/miniredis"
	// "github.com/go-redis/redis"
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
	mockRedis := mockstore.NewMockRedisStorage(ctrl)

	mockService := NewAuthService(mockStorage, mockRedis)

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

	sess := &types.Session{
		UserID: reqAcc.ID,
	}
	token := "refresh-token"
	mockRedis.EXPECT().CreateSession(context.Background(), gomock.Eq(sess), 86400).Return(token, nil)


	account, err := mockService.CreateAccount(context.Background(), req)
	require.NoError(t, err)
	require.Nil(t, err)
	require.Equal(t, accountToProto(acc), account.Account)
}

func Test_UpdateAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage, nil)

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

	// mr, err := miniredis.Run()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// client := redis.NewClient(&redis.Options{
	// 	Addr: mr.Addr(),
	// })

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockRedis := mockstore.NewMockRedisStorage(ctrl)

	mockService := NewAuthService(mockStorage, mockRedis)

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

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage, nil)

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

func Test_UpdateBalance(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage, nil)

	req := &authpb.UpdateBalanceRequest{
		Id:           uuid.New().String(),
		Balance:      50,
		BlockedMoney: 50,
	}
	uid, _ := uuid.Parse(req.Id)
	acc := &types.Account{
		ID:               uid,
		FirstName:        "Pasha1",
		LastName:         "Volkov",
		CardNumber:       "4444444444444443",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          50,
		BlockedMoney:     50,
		CreatedAt:        time.Now(),
	}
	mockStorage.EXPECT().SaveBalance(context.Background(), req).Return(acc, nil).AnyTimes()

	account, err := mockService.UpdateBalance(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, account.Balance, acc.Balance)
	require.Equal(t, account.BlockedMoney, acc.BlockedMoney)
}

func Test_GetAccountByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockService := NewAuthService(mockStorage, nil)

	req := &authpb.GetIDRequest{
		Id: uuid.New().String(),
	}

	uid, _ := uuid.Parse(req.Id)
	acc := &types.Account{
		ID:               uid,
		FirstName:        "Pasha1",
		LastName:         "Volkov",
		CardNumber:       "4444444444444443",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          50,
		BlockedMoney:     50,
		CreatedAt:        time.Now(),
	}

	mockStorage.EXPECT().GetAccountByID(context.Background(), req).Return(acc, nil).AnyTimes()

	account, err := mockService.GetAccountByID(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, account.Id, acc.ID.String())
	require.Equal(t, account.Balance, acc.Balance)
	require.Equal(t, account.BlockedMoney, acc.BlockedMoney)
}

func Test_GetAccount(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mockstore.NewMockStorage(ctrl)
	service := NewAuthService(storage, nil)
	acc1 := &types.Account{
		ID:               uuid.New(),
		FirstName:        "Pasha",
		LastName:         "Volkov",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          0,
		BlockedMoney:     0,
		Statement:        []string{},
		CreatedAt:        time.Now(),
	}
	acc2 := &types.Account{
		ID:               uuid.New(),
		FirstName:        "Pasha",
		LastName:         "Volkov",
		CardNumber:       "4444444444444443",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          0,
		BlockedMoney:     0,
		Statement:        []string{},
		CreatedAt:        time.Now(),
	}
	accs := []*types.Account{acc1, acc2}
	streamServer := mockclient.NewMockAuthService_GetAccountServer(ctrl)

	streamServer.EXPECT().Context().Return(context.Background()).AnyTimes()
	storage.EXPECT().GetAccount(context.Background()).Return(accs, nil).AnyTimes()
	streamServer.EXPECT().Send(gomock.Any()).DoAndReturn(
		func(req *authpb.Account) error {
			if req.Id == acc1.ID.String() {
				return nil
			}
			if req.Id == acc2.ID.String() {
				return nil
			}
			return fmt.Errorf("not found")
		},
	).AnyTimes()

	err := service.GetAccount(&authpb.GetRequest{}, streamServer)
	require.NoError(t, err)
	require.Nil(t, err)
}

func Test_CreateStatement(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mockstore.NewMockStorage(ctrl)
	service := NewAuthService(storage, nil)

	req1 := &authpb.StatementRequest{
		AccountId: uuid.New().String(),
		PaymentId: uuid.New().String(),
	}

	req2 := &authpb.StatementRequest{
		AccountId: uuid.New().String(),
		PaymentId: uuid.New().String(),
	}

	streamServer := mockclient.NewMockAuthService_CreateStatementServer(ctrl)

	streamServer.EXPECT().Context().Return(context.Background()).AnyTimes()

	streamServer.EXPECT().Recv().Return(req1, nil).AnyTimes()
	storage.EXPECT().UpdateStatement(context.Background(), req1).Return([]string{req1.PaymentId}, nil).AnyTimes()
	streamServer.EXPECT().Send(&authpb.StatementResponse{}).Return(nil).AnyTimes()

	streamServer.EXPECT().Recv().Return(req2, nil).AnyTimes()
	storage.EXPECT().UpdateStatement(context.Background(), req2).Return([]string{req2.PaymentId}, nil).AnyTimes()
	streamServer.EXPECT().Send(&authpb.StatementResponse{}).Return(nil).AnyTimes()

	err := service.CreateStatement(streamServer)
	require.NoError(t, err)
	require.Nil(t, err)
}

func Test_GetStatement(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mockstore.NewMockStorage(ctrl)
	service := NewAuthService(storage, nil)

	req := &authpb.StatementGet{
		AccountId: uuid.New().String(),
	}

	uid, _ := uuid.Parse(req.AccountId)
	st1 := &authpb.Statement{
		PaymentId: uuid.New().String(),
	}
	st2 := &authpb.Statement{
		PaymentId: uuid.New().String(),
	}
	acc := &types.Account{
		ID:               uid,
		FirstName:        "Pasha1",
		LastName:         "Volkov1",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "123",
		Balance:          0,
		BlockedMoney:     0,
		Statement:        []string{st1.PaymentId, st2.PaymentId},
		CreatedAt:        time.Now(),
	}

	streamServer := mockclient.NewMockAuthService_GetStatementServer(ctrl)

	storage.EXPECT().GetAccountByID(context.Background(), &authpb.GetIDRequest{
		Id: req.AccountId,
	}).Return(acc, nil).AnyTimes()
	streamServer.EXPECT().Context().Return(context.Background()).AnyTimes()
	streamServer.EXPECT().Send(gomock.Any()).DoAndReturn(
		func(req *authpb.Statement) error {
			if req.PaymentId == st1.PaymentId {
				return nil
			}
			if req.PaymentId == st2.PaymentId {
				return nil
			}
			return fmt.Errorf("not found")
		},
	).AnyTimes()

	err := service.GetStatement(req, streamServer)
	require.NoError(t, err)
	require.Nil(t, err)
}

func Test_RefreshTokens(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockRedis := mockstore.NewMockRedisStorage(ctrl)

	mockService := NewAuthService(mockStorage, mockRedis)

	req := &authpb.RefreshRequest{
		RefreshToken: "cookieValue",
	}

	uid := uuid.New()
	mockRedis.EXPECT().GetUserID(context.Background(), req.RefreshToken).Return(uid, nil).AnyTimes()

	account := &types.Account{
		ID:               uid,
		FirstName:        "Pasha1",
		LastName:         "volkov1",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "924",
		Balance:          0,
		BlockedMoney:     0,
		Statement:        []string{},
		CreatedAt:        time.Now(),
	}

	mockStorage.EXPECT().GetAccountByID(context.Background(), gomock.Any()).Return(account, nil).AnyTimes()

	token := "refresh-token"
	sess := &types.Session{
		UserID: uid,
	}
	mockRedis.EXPECT().CreateSession(context.Background(), gomock.Eq(sess), 86400).Return(token, nil).AnyTimes()


	tokens, err := mockService.RefreshTokens(context.Background(), req)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, tokens)
}

func Test_SignIn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockRedis := mockstore.NewMockRedisStorage(ctrl)

	mockService := NewAuthService(mockStorage, mockRedis)
	uid := uuid.New()

	req := &authpb.LoginRequest{
		Id: uid.String(),
	}

	account := &types.Account{
		ID:               uid,
		FirstName:        "Pasha1",
		LastName:         "volkov1",
		CardNumber:       "4444444444444444",
		CardExpiryMonth:  "12",
		CardExpiryYear:   "24",
		CardSecurityCode: "924",
		Balance:          0,
		BlockedMoney:     0,
		Statement:        []string{},
		CreatedAt:        time.Now(),
	}

	mockStorage.EXPECT().GetAccountByID(context.Background(), gomock.Any()).Return(account, nil).AnyTimes()
	sess := &types.Session{
		UserID: uid,
	}
	token := "refresh-token"
	mockRedis.EXPECT().CreateSession(context.Background(), gomock.Eq(sess), 86400).Return(token, nil).AnyTimes()

	accWithTokens, err := mockService.SignIn(context.Background(), req)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, accWithTokens)
}

func Test_SignOut(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mockStorage := mockstore.NewMockStorage(ctrl)
	mockRedis := mockstore.NewMockRedisStorage(ctrl)

	mockService := NewAuthService(mockStorage, mockRedis)

	cookieValue := "cookieValue"

	mockRedis.EXPECT().DeleteSession(context.Background(), gomock.Eq(cookieValue)).Return(nil).AnyTimes()

	message, err := mockService.SignOut(context.Background(), &authpb.QuitRequest{
		RefreshToken: cookieValue,
	})
	require.NoError(t, err)
	require.Nil(t, err)
	require.Equal(t, message.Message, "sign-out")
}