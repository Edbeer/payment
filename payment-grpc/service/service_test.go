package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mock_proto "github.com/Edbeer/auth/client/mock"
	authpb "github.com/Edbeer/auth/proto"

	paymentpb "github.com/Edbeer/payment-grpc/proto"
	mockpay "github.com/Edbeer/payment-grpc/service/mock"

	"github.com/Edbeer/payment-grpc/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_CreatePayment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)
		req := &paymentpb.CreateRequest{
			Merchant:         uuid.New().String(),
			Customer:         uuid.New().String(),
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Currency:         "rub",
			Amount:           50,
		}

		reqIDC := &authpb.GetIDRequest{
			Id: req.Customer,
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: req.Merchant,
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		mock.ExpectBegin()
		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		sts := []*authpb.StatementRequest{}

		customer.Balance = customer.Balance - req.Amount
		customer.BlockedMoney = customer.BlockedMoney + req.Amount
		merchant.BlockedMoney = merchant.BlockedMoney + req.Amount

		clientAuth.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, req *authpb.UpdateBalanceRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
				if req.Id == customer.Id {
					return customer, nil
				}
				if req.Id == merchant.Id {
					return merchant, nil
				}
				return nil, fmt.Errorf("not found")
			},
		).AnyTimes()

		payment := types.CreateAuthPayment(req, customer, merchant, "Approved")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(payment, nil).AnyTimes()

		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: payment.PaymentId.String(),
		})
		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: payment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).DoAndReturn(
			func(req *authpb.StatementRequest) error {
				if req.AccountId == customer.Id {
					return nil
				}
				if req.AccountId == merchant.Id {
					return nil
				}
				return fmt.Errorf("not found")
			},
		).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()
		st, err := servicePay.CreatePayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, payment.PaymentId.String())
		require.Equal(t, st.Status, "Approved")
	})

	t.Run("Wrong Payment request", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)
		req := &paymentpb.CreateRequest{
			Merchant:         uuid.New().String(),
			Customer:         uuid.New().String(),
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Currency:         "rub",
			Amount:           50,
		}

		reqIDC := &authpb.GetIDRequest{
			Id: req.Customer,
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444423",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: req.Merchant,
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		require.NotEqual(t, customer.CardNumber, req.CardNumber)

		mock.ExpectBegin()

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		payment := types.CreateAuthPayment(req, customer, merchant, "wrong payment request")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(payment, nil).AnyTimes()

		sts := []*authpb.StatementRequest{}
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: payment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(&authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: payment.PaymentId.String(),
		}).Return(nil).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()
		st, err := servicePay.CreatePayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, payment.PaymentId.String())
		require.Equal(t, st.Status, "wrong payment request")
	})

	t.Run("Insufficient funds", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.CreateRequest{
			Merchant:         uuid.New().String(),
			Customer:         uuid.New().String(),
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Currency:         "rub",
			Amount:           50,
		}

		reqIDC := &authpb.GetIDRequest{
			Id: req.Customer,
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          40,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: req.Merchant,
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		mock.ExpectBegin()

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		payment := types.CreateAuthPayment(req, customer, merchant, "Insufficient funds")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(payment, nil).AnyTimes()

		sts := []*authpb.StatementRequest{}
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: payment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(&authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: payment.PaymentId.String(),
		}).Return(nil).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()
		st, err := servicePay.CreatePayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, payment.PaymentId.String())
		require.Equal(t, st.Status, "Insufficient funds")
	})
}

func Test_CapturePayment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successful payment", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
			Amount:    50,
		}
		pid, _ := uuid.Parse(req.PaymentId)

		refPayment := &types.Payment{
			PaymentId:       pid,
			Merchant:        uuid.New(),
			Customer:        uuid.New(),
			CardNumber:      "4444444444444444",
			CardExpiryMonth: "12",
			CardExpiryYear:  "24",
			Currency:        "rub",
			Operation:       "Authorization",
			Status:          "Approved",
			Amount:          50,
			CreatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		storagePay.EXPECT().GetPaymentByID(context.Background(), req).Return(refPayment, nil).AnyTimes()

		reqIDC := &authpb.GetIDRequest{
			Id: refPayment.Customer.String(),
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: refPayment.Merchant.String(),
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		newPayment := types.CreateCompletePayment(req, refPayment, "Successful payment")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(newPayment, nil).AnyTimes()

		// update balance
		customer.BlockedMoney = customer.BlockedMoney - req.Amount

		merchant.Balance = merchant.Balance + req.Amount
		merchant.BlockedMoney = merchant.BlockedMoney - req.Amount

		clientAuth.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, req *authpb.UpdateBalanceRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
				if req.Id == customer.Id {
					return customer, nil
				}
				if req.Id == merchant.Id {
					return merchant, nil
				}
				return nil, fmt.Errorf("not found")
			},
		).AnyTimes()

		sts := []*authpb.StatementRequest{}

		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: newPayment.PaymentId.String(),
		})
		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: newPayment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).DoAndReturn(
			func(req *authpb.StatementRequest) error {
				if req.AccountId == customer.Id {
					return nil
				}
				if req.AccountId == merchant.Id {
					return nil
				}
				return fmt.Errorf("not found")
			},
		).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()
		st, err := servicePay.CapturePayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, newPayment.PaymentId.String())
		require.Equal(t, st.Status, "Successful payment")
	})

	t.Run("Invalid amount", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
			Amount:    60,
		}
		pid, _ := uuid.Parse(req.PaymentId)

		refPayment := &types.Payment{
			PaymentId:       pid,
			Merchant:        uuid.New(),
			Customer:        uuid.New(),
			CardNumber:      "4444444444444444",
			CardExpiryMonth: "12",
			CardExpiryYear:  "24",
			Currency:        "rub",
			Operation:       "Authorization",
			Status:          "Approved",
			Amount:          50,
			CreatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		storagePay.EXPECT().GetPaymentByID(context.Background(), req).Return(refPayment, nil).AnyTimes()

		reqIDC := &authpb.GetIDRequest{
			Id: refPayment.Customer.String(),
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: refPayment.Merchant.String(),
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		newPayment := types.CreateCompletePayment(req, refPayment, "Invalid amount")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(newPayment, nil).AnyTimes()
		sts := []*authpb.StatementRequest{}
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: newPayment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()

		st, err := servicePay.CapturePayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, newPayment.PaymentId.String())
		require.Equal(t, st.Status, "Invalid amount")
	})
}

func Test_RefundPayment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successful refund", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
			Amount:    50,
		}
		pid, _ := uuid.Parse(req.PaymentId)

		refPayment := &types.Payment{
			PaymentId:       pid,
			Merchant:        uuid.New(),
			Customer:        uuid.New(),
			CardNumber:      "4444444444444444",
			CardExpiryMonth: "12",
			CardExpiryYear:  "24",
			Currency:        "rub",
			Operation:       "Capture",
			Status:          "Successful payment",
			Amount:          50,
			CreatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		storagePay.EXPECT().GetPaymentByID(context.Background(), req).Return(refPayment, nil).AnyTimes()

		reqIDC := &authpb.GetIDRequest{
			Id: refPayment.Customer.String(),
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: refPayment.Merchant.String(),
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		newPayment := types.CreateCompletePayment(req, refPayment, "Successful refund")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(newPayment, nil).AnyTimes()

		// update balance
		customer.Balance = customer.Balance + req.Amount

		merchant.Balance = merchant.Balance - req.Amount

		clientAuth.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, req *authpb.UpdateBalanceRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
				if req.Id == customer.Id {
					return customer, nil
				}
				if req.Id == merchant.Id {
					return merchant, nil
				}
				return nil, fmt.Errorf("not found")
			},
		).AnyTimes()

		sts := []*authpb.StatementRequest{}

		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: newPayment.PaymentId.String(),
		})
		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: newPayment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).DoAndReturn(
			func(req *authpb.StatementRequest) error {
				if req.AccountId == customer.Id {
					return nil
				}
				if req.AccountId == merchant.Id {
					return nil
				}
				return fmt.Errorf("not found")
			},
		).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()

		st, err := servicePay.RefundPayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, newPayment.PaymentId.String())
		require.Equal(t, st.Status, "Successful refund")
	})

	t.Run("Ivalid amount", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
			Amount:    60,
		}
		pid, _ := uuid.Parse(req.PaymentId)

		refPayment := &types.Payment{
			PaymentId:       pid,
			Merchant:        uuid.New(),
			Customer:        uuid.New(),
			CardNumber:      "4444444444444444",
			CardExpiryMonth: "12",
			CardExpiryYear:  "24",
			Currency:        "rub",
			Operation:       "Capture",
			Status:          "Successful payment",
			Amount:          50,
			CreatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		storagePay.EXPECT().GetPaymentByID(context.Background(), req).Return(refPayment, nil).AnyTimes()

		reqIDC := &authpb.GetIDRequest{
			Id: refPayment.Customer.String(),
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: refPayment.Merchant.String(),
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		newPayment := types.CreateCompletePayment(req, refPayment, "Invalid amount")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(newPayment, nil).AnyTimes()
		sts := []*authpb.StatementRequest{}
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: newPayment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()
		mock.ExpectCommit()

		st, err := servicePay.RefundPayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, newPayment.PaymentId.String())
		require.Equal(t, st.Status, "Invalid amount")
	})
}

func Test_CancelPayment(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successful cancel", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
			Amount:    50,
		}
		pid, _ := uuid.Parse(req.PaymentId)

		refPayment := &types.Payment{
			PaymentId:       pid,
			Merchant:        uuid.New(),
			Customer:        uuid.New(),
			CardNumber:      "4444444444444444",
			CardExpiryMonth: "12",
			CardExpiryYear:  "24",
			Currency:        "rub",
			Operation:       "Authorization",
			Status:          "Approved",
			Amount:          50,
			CreatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		storagePay.EXPECT().GetPaymentByID(context.Background(), req).Return(refPayment, nil).AnyTimes()

		reqIDC := &authpb.GetIDRequest{
			Id: refPayment.Customer.String(),
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: refPayment.Merchant.String(),
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		newPayment := types.CreateCompletePayment(req, refPayment, "Successful cancel")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(newPayment, nil).AnyTimes()

		// update balance
		customer.Balance = customer.Balance + req.Amount
		customer.BlockedMoney = customer.BlockedMoney - req.Amount

		merchant.BlockedMoney = merchant.BlockedMoney - req.Amount

		clientAuth.EXPECT().UpdateBalance(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, req *authpb.UpdateBalanceRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
				if req.Id == customer.Id {
					return customer, nil
				}
				if req.Id == merchant.Id {
					return merchant, nil
				}
				return nil, fmt.Errorf("not found")
			},
		).AnyTimes()

		sts := []*authpb.StatementRequest{}

		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: newPayment.PaymentId.String(),
		})
		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: newPayment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).DoAndReturn(
			func(req *authpb.StatementRequest) error {
				if req.AccountId == customer.Id {
					return nil
				}
				if req.AccountId == merchant.Id {
					return nil
				}
				return fmt.Errorf("not found")
			},
		).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()

		mock.ExpectCommit()

		st, err := servicePay.CancelPayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, newPayment.PaymentId.String())
		require.Equal(t, st.Status, "Successful cancel")
	})

	t.Run("Ivalid amount", func(t *testing.T) {
		storagePay := mockpay.NewMockStorage(ctrl)
		clientAuth := mock_proto.NewMockAuthServiceClient(ctrl)

		servicePay := NewPaymentService(storagePay, clientAuth, db)

		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
			Amount:    60,
		}
		pid, _ := uuid.Parse(req.PaymentId)

		refPayment := &types.Payment{
			PaymentId:       pid,
			Merchant:        uuid.New(),
			Customer:        uuid.New(),
			CardNumber:      "4444444444444444",
			CardExpiryMonth: "12",
			CardExpiryYear:  "24",
			Currency:        "rub",
			Operation:       "Authorization",
			Status:          "Approved",
			Amount:          50,
			CreatedAt:       time.Now(),
		}

		mock.ExpectBegin()
		storagePay.EXPECT().GetPaymentByID(context.Background(), req).Return(refPayment, nil).AnyTimes()

		reqIDC := &authpb.GetIDRequest{
			Id: refPayment.Customer.String(),
		}
		customer := &authpb.Account{
			Id:               reqIDC.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "999",
			Balance:          100,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		reqIDM := &authpb.GetIDRequest{
			Id: refPayment.Merchant.String(),
		}
		merchant := &authpb.Account{
			Id:               reqIDM.Id,
			FirstName:        "Pasha1",
			LastName:         "Volkov",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
			Balance:          50,
			BlockedMoney:     50,
			Statement:        []string{},
			CreatedAt:        timestamppb.Now(),
		}

		clientAuth.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, req *authpb.GetIDRequest, opts ...grpc.CallOption) (*authpb.Account, error) {
					if req.Id == reqIDC.Id {
						return customer, nil
					}
					if req.Id == reqIDM.Id {
						return merchant, nil
					}
					return nil, fmt.Errorf("not found")
				},
			).AnyTimes()

		newPayment := types.CreateCompletePayment(req, refPayment, "Invalid amount")
		storagePay.EXPECT().SavePayment(context.Background(), gomock.Any(), gomock.Any()).Return(newPayment, nil).AnyTimes()
		sts := []*authpb.StatementRequest{}
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: newPayment.PaymentId.String(),
		})

		streamSts := mock_proto.NewMockAuthService_CreateStatementClient(ctrl)
		clientAuth.EXPECT().CreateStatement(context.Background()).Return(streamSts, nil).AnyTimes()
		streamSts.EXPECT().Send(gomock.Any()).Return(nil).AnyTimes()
		streamSts.EXPECT().Recv().Return(&authpb.StatementResponse{}, nil).AnyTimes()
		streamSts.EXPECT().CloseSend().Return(nil).AnyTimes()

		mock.ExpectCommit()

		st, err := servicePay.CancelPayment(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, st)
		require.Equal(t, st.PaymentId, newPayment.PaymentId.String())
		require.Equal(t, st.Status, "Invalid amount")
	})
}
