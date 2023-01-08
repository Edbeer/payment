//go:generate mockgen -source service.go -destination mock/storage_mock.go -package mock
package service

import (
	"context"
	"database/sql"

	authpb "github.com/Edbeer/auth-grpc/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/types"
)

type Storage interface {
	SavePayment(ctx context.Context, tx *sql.Tx, payment *types.Payment) (*types.Payment, error)
}

type PaymentService struct {
	paymentpb.UnimplementedPaymentServiceServer
	client  authpb.AuthServiceClient
	storage Storage
	db      *sql.DB
}

func NewPaymentService(storage Storage, client authpb.AuthServiceClient, db *sql.DB) *PaymentService {
	return &PaymentService{storage: storage, client: client, db: db}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *paymentpb.CreateRequest) (*paymentpb.Statement, error) {
	// get customer
	customer, err := getCustomerByID(ctx, s.client, req)
	if err != nil {
		return nil, err
	}
	// get merchant
	merchant, err := getMerchantByID(ctx, s.client, req)
	if err != nil {
		return nil, err
	}
	sts := []*authpb.StatementRequest{}
	// check payment request
	if req.CardNumber != customer.CardNumber ||
		req.CardExpiryMonth != customer.CardExpiryMonth ||
		req.CardExpiryYear != customer.CardExpiryYear ||
		req.CardSecurityCode != customer.CardSecurityCode {
		// Begin Transaction
		tx, err := s.db.BeginTx(ctx, nil)
		defer tx.Rollback()
		if err != nil {
			return nil, err
		}
		// create payment
		payment := types.CreateAuthPayment(req, customer, merchant, "wrong payment request")
		savedPayment, err := s.storage.SavePayment(ctx, tx, payment)
		if err != nil {
			return nil, err
		}
		// send statements to auth service
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: savedPayment.PaymentId.String(),
		})
		if err := createStatement(ctx, s.client, sts); err != nil {
			tx.Rollback()
			return nil, err
		}
		// // send statements to auth service
		// _, err = s.client.CreateStatement(ctx, &authpb.StatementRequest{
		// 	AccountId: merchant.Id,
		// 	PaymentId: savedPayment.PaymentId.String(),
		// })
		// if err != nil {
		// 	tx.Rollback()
		// 	return nil, err
		// }
		// Commit transaction
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &paymentpb.Statement{
			PaymentId: savedPayment.PaymentId.String(),
			Status:    savedPayment.Status,
		}, nil
	}
	// consume customer balance
	// balance < req amount
	if customer.Balance < req.Amount {
		// Begin Transaction
		tx, err := s.db.BeginTx(ctx, nil)
		defer tx.Rollback()
		if err != nil {
			return nil, err
		}
		// create payment
		payment := types.CreateAuthPayment(req, customer, merchant, "Insufficient funds")
		savedPayment, err := s.storage.SavePayment(ctx, tx, payment)
		if err != nil {
			return nil, err
		}
		// save statement for merchant
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: savedPayment.PaymentId.String(),
		})
		if err := createStatement(ctx, s.client, sts); err != nil {
			tx.Rollback()
			return nil, err
		}
		// // save statement for merchant
		// _, err = s.client.CreateStatement(ctx, &authpb.StatementRequest{
		// 	AccountId: merchant.Id,
		// 	PaymentId: savedPayment.PaymentId.String(),
		// })
		// if err != nil {
		// 	tx.Rollback()
		// 	return nil, err
		// }
		// Commit transaction
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &paymentpb.Statement{
			PaymentId: savedPayment.PaymentId.String(),
			Status:    savedPayment.Status,
		}, nil
	}
	// balance > req amount
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}

	// customer acc new balance
	customer, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
		Id:           customer.Id,
		Balance:      customer.Balance - req.Amount,
		BlockedMoney: customer.BlockedMoney + req.Amount,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// merchant acc new balance
	merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
		Id:           merchant.Id,
		Balance:      merchant.Balance,
		BlockedMoney: merchant.BlockedMoney + req.Amount,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// create new payment
	payment := types.CreateAuthPayment(req, customer, merchant, "Approved")
	savedPayment, err := s.storage.SavePayment(ctx, tx, payment)
	if err != nil {
		return nil, err
	}
	// // save statement for customer
	// _, err = s.client.CreateStatement(ctx, &authpb.StatementRequest{
	// 	AccountId: customer.Id,
	// 	PaymentId: savedPayment.PaymentId.String(),
	// })
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }
	// save statement for customer
	sts = append(sts, &authpb.StatementRequest{
		AccountId: customer.Id,
		PaymentId: savedPayment.PaymentId.String(),
	})
	if err := createStatement(ctx, s.client, sts); err != nil {
		tx.Rollback()
		return nil, err
	}
	// save statement for merchant
	sts = append(sts, &authpb.StatementRequest{
		AccountId: merchant.Id,
		PaymentId: savedPayment.PaymentId.String(),
	})
	if err := createStatement(ctx, s.client, sts); err != nil {
		tx.Rollback()
		return nil, err
	}
	// // save statement for merchant
	// _, err = s.client.CreateStatement(ctx, &authpb.StatementRequest{
	// 	AccountId: merchant.Id,
	// 	PaymentId: savedPayment.PaymentId.String(),
	// })
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &paymentpb.Statement{
		PaymentId: savedPayment.PaymentId.String(),
		Status:    savedPayment.Status,
	}, nil
}

// get customer account
func getCustomerByID(ctx context.Context, client authpb.AuthServiceClient, req *paymentpb.CreateRequest) (*authpb.Account, error) {
	account, err := client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: req.Customer,
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

// get merchant  account
func getMerchantByID(ctx context.Context, client authpb.AuthServiceClient, req *paymentpb.CreateRequest) (*authpb.Account, error) {
	account, err := client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: req.Merchant,
	})
	if err != nil {
		return nil, err
	}
	return account, nil
}

func createStatement(ctx context.Context, client authpb.AuthServiceClient, sts []*authpb.StatementRequest) error {
	stream, err := client.CreateStatement(ctx)
	if err != nil {
		return err
	}
	for _, statement := range sts {
		if err := stream.Send(statement); err != nil {
			return err
		}
		_,  err := stream.Recv(); 
		if err != nil {
			return err
		}
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	return nil
}

// func getAccountByID(client authpb.AuthServiceClient, req *paymentpb.Payment) ([]*authpb.Account, error) {
// 	stream, err := client.GetAccountByID(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	accsID := []*authpb.GetIDRequest{
// 		{
// 			Id: req.Customer,
// 		},
// 		{
// 			Id: req.Merchant,
// 		},
// 	}
// 	for _, acc := range accsID {
// 		if err := stream.Send(acc); err != nil {
// 			return nil, err
// 		}
// 	}
// 	if err := stream.CloseSend(); err != nil {
// 		return nil, err
// 	}
// 	accounts := []*authpb.Account{}
// 	for {
// 		account, err := stream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		accounts = append(accounts, account)
// 	}
// 	return accounts, nil
// }

func (s *PaymentService) CapturePayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}

func (s *PaymentService) CancelPayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, req *paymentpb.Payment) (*paymentpb.Statement, error) {
	return nil, nil
}
