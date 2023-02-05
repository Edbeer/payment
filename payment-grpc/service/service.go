//go:generate mockgen -source service.go -destination mock/storage_mock.go -package mock
package service

import (
	"context"
	"database/sql"

	authpb "github.com/Edbeer/auth/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/types"
)

type Storage interface {
	SavePayment(ctx context.Context, payment *types.Payment, tx *sql.Tx) (*types.Payment, error)
	GetPaymentByID(ctx context.Context, req *paymentpb.PaidRequest) (*types.Payment, error)
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
	// Begin Tx
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// get customer
	customer, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: req.Customer,
	})
	if err != nil {
		return nil, err
	}
	// get merchant
	merchant, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: req.Merchant,
	})
	if err != nil {
		return nil, err
	}
	// statements
	sts := []*authpb.StatementRequest{}
	// check payment request
	if req.CardNumber != customer.CardNumber ||
		req.CardExpiryMonth != customer.CardExpiryMonth ||
		req.CardExpiryYear != customer.CardExpiryYear ||
		req.CardSecurityCode != customer.CardSecurityCode {
		// create payment
		payment := types.CreateAuthPayment(req, customer, merchant, "wrong payment request")
		savedPayment, err := s.storage.SavePayment(ctx, payment, tx)
		if err != nil {
			return nil, err
		}
		// send statements to auth service
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: savedPayment.PaymentId.String(),
		})
		if err := createStatement(ctx, s.client, sts); err != nil {
			return nil, err
		}
		// commit tx
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
		// create payment
		payment := types.CreateAuthPayment(req, customer, merchant, "Insufficient funds")
		savedPayment, err := s.storage.SavePayment(ctx, payment, tx)
		if err != nil {
			return nil, err
		}
		// save statement for merchant
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: savedPayment.PaymentId.String(),
		})
		// send statements to auth service
		if err := createStatement(ctx, s.client, sts); err != nil {
			return nil, err
		}
		// commit tx
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &paymentpb.Statement{
			PaymentId: savedPayment.PaymentId.String(),
			Status:    savedPayment.Status,
		}, nil
	}
	// balance > req amount
	// customer acc new balance
	customer.Balance = customer.Balance - req.Amount
	customer.BlockedMoney = customer.BlockedMoney + req.Amount
	customer, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
		Id:           customer.Id,
		Balance:      customer.Balance,
		BlockedMoney: customer.BlockedMoney,
	})
	if err != nil {
		return nil, err
	}
	// merchant acc new balance
	merchant.BlockedMoney = merchant.BlockedMoney + req.Amount
	merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
		Id:           merchant.Id,
		Balance:      merchant.Balance,
		BlockedMoney: merchant.BlockedMoney,
	})
	if err != nil {
		customer.Balance = customer.Balance + req.Amount
		customer.BlockedMoney = customer.BlockedMoney - req.Amount
		customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           customer.Id,
			Balance:      customer.Balance,
			BlockedMoney: customer.BlockedMoney,
		})

		return nil, err
	}
	// create new payment
	payment := types.CreateAuthPayment(req, customer, merchant, "Approved")
	savedPayment, err := s.storage.SavePayment(ctx, payment, tx)
	if err != nil {
		customer.Balance = customer.Balance + req.Amount
		customer.BlockedMoney = customer.BlockedMoney - req.Amount
		customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           customer.Id,
			Balance:      customer.Balance,
			BlockedMoney: customer.BlockedMoney,
		})

		merchant.BlockedMoney = merchant.BlockedMoney - req.Amount
		merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           merchant.Id,
			Balance:      merchant.Balance,
			BlockedMoney: merchant.BlockedMoney,
		})
		return nil, err
	}
	// save statement for customer
	sts = append(sts, &authpb.StatementRequest{
		AccountId: customer.Id,
		PaymentId: savedPayment.PaymentId.String(),
	})
	// save statement for merchant
	sts = append(sts, &authpb.StatementRequest{
		AccountId: merchant.Id,
		PaymentId: savedPayment.PaymentId.String(),
	})
	// send statements to auth service
	if err := createStatement(ctx, s.client, sts); err != nil {
		customer.Balance = customer.Balance + req.Amount
		customer.BlockedMoney = customer.BlockedMoney - req.Amount
		customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           customer.Id,
			Balance:      customer.Balance,
			BlockedMoney: customer.BlockedMoney,
		})

		merchant.BlockedMoney = merchant.BlockedMoney - req.Amount
		merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           merchant.Id,
			Balance:      merchant.Balance,
			BlockedMoney: merchant.BlockedMoney,
		})
		return nil, err
	}
	// commit tx
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &paymentpb.Statement{
		PaymentId: savedPayment.PaymentId.String(),
		Status:    savedPayment.Status,
	}, nil
}

func (s *PaymentService) CapturePayment(ctx context.Context, req *paymentpb.PaidRequest) (*paymentpb.Statement, error) {
	// Begin Tx
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// Get referenced payment
	refPayment, err := s.storage.GetPaymentByID(ctx, req)
	if err != nil {
		return nil, err
	}
	// Get customer
	customer, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: refPayment.Customer.String(),
	})
	// Get merchant
	merchant, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: refPayment.Merchant.String(),
	})
	// statements
	sts := []*authpb.StatementRequest{}
	if refPayment.Operation == "Authorization" && refPayment.Status == "Approved" {
		// Invalid amount
		if refPayment.Amount < req.Amount {
			refPayment.Operation = "Capture"
			completedPayment := types.CreateCompletePayment(req, refPayment, "Invalid amount")
			invalidPayment, err := s.storage.SavePayment(ctx, completedPayment, tx)
			if err != nil {
				return nil, err
			}
			// send statements to auth service
			sts = append(sts, &authpb.StatementRequest{
				AccountId: merchant.Id,
				PaymentId: completedPayment.PaymentId.String(),
			})
			if err := createStatement(ctx, s.client, sts); err != nil {
				return nil, err
			}
			// commit tx
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return &paymentpb.Statement{
				PaymentId: invalidPayment.PaymentId.String(),
				Status:    invalidPayment.Status,
			}, nil
		}
		// Successful payment
		// update customer balance and append new statement
		customer.BlockedMoney = customer.BlockedMoney - req.Amount

		customer, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           customer.Id,
			BlockedMoney: customer.BlockedMoney,
		})
		if err != nil {
			return nil, err
		}
		// update new merchant balance and append new statement
		merchant.Balance = merchant.Balance + req.Amount
		merchant.BlockedMoney = merchant.BlockedMoney - req.Amount

		merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           merchant.Id,
			Balance:      merchant.Balance,
			BlockedMoney: merchant.BlockedMoney,
		})
		if err != nil {
			customer.BlockedMoney = customer.BlockedMoney + req.Amount

			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           customer.Id,
				BlockedMoney: customer.BlockedMoney,
			})
			return nil, err
		}
		// make complete payment
		completedPayment := types.CreateCompletePayment(req, refPayment, "Successful payment")
		completedPayment.Operation = "Capture"
		completedPayment, err = s.storage.SavePayment(ctx, completedPayment, tx)
		if err != nil {
			customer.BlockedMoney = customer.BlockedMoney + req.Amount

			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           customer.Id,
				BlockedMoney: customer.BlockedMoney,
			})

			merchant.Balance = merchant.Balance - req.Amount
			merchant.BlockedMoney = merchant.BlockedMoney + req.Amount

			merchant, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           merchant.Id,
				Balance:      merchant.Balance,
				BlockedMoney: merchant.BlockedMoney,
			})
			return nil, err
		}
		// save statement for customer
		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: completedPayment.PaymentId.String(),
		})
		// save statement for merchant
		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: completedPayment.PaymentId.String(),
		})
		// send statements to auth service
		if err := createStatement(ctx, s.client, sts); err != nil {
			customer.BlockedMoney = customer.BlockedMoney + req.Amount

			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           customer.Id,
				BlockedMoney: customer.BlockedMoney,
			})

			merchant.Balance = merchant.Balance - req.Amount
			merchant.BlockedMoney = merchant.BlockedMoney + req.Amount

			merchant, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           merchant.Id,
				Balance:      merchant.Balance,
				BlockedMoney: merchant.BlockedMoney,
			})
			return nil, err
		}
		// commit tx
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &paymentpb.Statement{
			PaymentId: completedPayment.PaymentId.String(),
			Status:    completedPayment.Status,
		}, nil
	}
	return &paymentpb.Statement{
		PaymentId: req.PaymentId,
		Status:    "Invalid transaction",
	}, nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, req *paymentpb.PaidRequest) (*paymentpb.Statement, error) {
	// Begin Tx
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// Get referenced payment
	refPayment, err := s.storage.GetPaymentByID(ctx, req)
	if err != nil {
		return nil, err
	}
	// Get customer
	customer, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: refPayment.Customer.String(),
	})
	// Get merchant
	merchant, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: refPayment.Merchant.String(),
	})
	// statements
	sts := []*authpb.StatementRequest{}
	if refPayment.Operation == "Capture" && refPayment.Status == "Successful payment" {
		// Invalid amount
		if refPayment.Amount < req.Amount {
			refPayment.Operation = "Refund"
			completedPayment := types.CreateCompletePayment(req, refPayment, "Invalid amount")
			invalidPayment, err := s.storage.SavePayment(ctx, completedPayment, tx)
			if err != nil {
				return nil, err
			}
			// send statements to auth service
			sts = append(sts, &authpb.StatementRequest{
				AccountId: merchant.Id,
				PaymentId: completedPayment.PaymentId.String(),
			})
			if err := createStatement(ctx, s.client, sts); err != nil {
				return nil, err
			}
			// commit tx
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return &paymentpb.Statement{
				PaymentId: invalidPayment.PaymentId.String(),
				Status:    invalidPayment.Status,
			}, nil
		}
		// Successful refund
		// update customer balance
		customer.Balance = customer.Balance + req.Amount

		customer, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:      customer.Id,
			Balance: customer.Balance,
		})
		if err != nil {
			return nil, err
		}
		// update new merchant balance
		merchant.Balance = merchant.Balance - req.Amount

		merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:      merchant.Id,
			Balance: merchant.Balance,
		})
		if err != nil {
			customer.Balance = customer.Balance - req.Amount

			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      customer.Id,
				Balance: customer.Balance,
			})
			return nil, err
		}
		// make complete refund
		refPayment.Operation = "Refund"
		completedPayment := types.CreateCompletePayment(req, refPayment, "Successful refund")
		completedPayment, err = s.storage.SavePayment(ctx, completedPayment, tx)
		if err != nil {
			customer.Balance = customer.Balance - req.Amount

			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      customer.Id,
				Balance: customer.Balance,
			})

			merchant.Balance = merchant.Balance + req.Amount

			merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      merchant.Id,
				Balance: merchant.Balance,
			})
			return nil, err
		}

		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: completedPayment.PaymentId.String(),
		})

		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: completedPayment.PaymentId.String(),
		})
		// send statements to auth service
		if err := createStatement(ctx, s.client, sts); err != nil {
			customer.Balance = customer.Balance - req.Amount

			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      customer.Id,
				Balance: customer.Balance,
			})

			merchant.Balance = merchant.Balance + req.Amount

			merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      merchant.Id,
				Balance: merchant.Balance,
			})
			return nil, err
		}
		// commit tx
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &paymentpb.Statement{
			PaymentId: completedPayment.PaymentId.String(),
			Status:    completedPayment.Status,
		}, nil
	}
	return &paymentpb.Statement{
		PaymentId: req.PaymentId,
		Status:    "Invalid transaction",
	}, nil
}

func (s *PaymentService) CancelPayment(ctx context.Context, req *paymentpb.PaidRequest) (*paymentpb.Statement, error) {
	// Begin Tx
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// Get referenced payment
	refPayment, err := s.storage.GetPaymentByID(ctx, req)
	if err != nil {
		return nil, err
	}
	// Get customer
	customer, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: refPayment.Customer.String(),
	})
	// Get merchant
	merchant, err := s.client.GetAccountByID(ctx, &authpb.GetIDRequest{
		Id: refPayment.Merchant.String(),
	})
	// statements
	sts := []*authpb.StatementRequest{}
	if refPayment.Operation == "Authorization" && refPayment.Status == "Approved" {
		// Invalid amount
		if refPayment.Amount < req.Amount {
			refPayment.Operation = "Cancel"
			completedPayment := types.CreateCompletePayment(req, refPayment, "Invalid amount")
			invalidPayment, err := s.storage.SavePayment(ctx, completedPayment, tx)
			if err != nil {
				return nil, err
			}
			// send statements to auth service
			sts = append(sts, &authpb.StatementRequest{
				AccountId: merchant.Id,
				PaymentId: completedPayment.PaymentId.String(),
			})
			if err := createStatement(ctx, s.client, sts); err != nil {
				return nil, err
			}
			// commit tx
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return &paymentpb.Statement{
				PaymentId: invalidPayment.PaymentId.String(),
				Status:    invalidPayment.Status,
			}, nil
		}
		// Successful cancel
		// update customer balance and append new statement
		customer.Balance = customer.Balance + req.Amount
		customer.BlockedMoney = customer.BlockedMoney - req.Amount
		customer, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:           customer.Id,
			Balance:      customer.Balance,
			BlockedMoney: customer.BlockedMoney,
		})
		if err != nil {
			return nil, err
		}

		// update new merchant balance and append new statement
		merchant.BlockedMoney = merchant.BlockedMoney - req.Amount

		merchant, err = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
			Id:      merchant.Id,
			Balance: merchant.BlockedMoney,
		})
		if err != nil {
			customer.Balance = customer.Balance - req.Amount
			customer.BlockedMoney = customer.BlockedMoney + req.Amount
			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           customer.Id,
				Balance:      customer.Balance,
				BlockedMoney: customer.BlockedMoney,
			})
			return nil, err
		}
		// make cancel
		refPayment.Operation = "Cancel"
		completedPayment := types.CreateCompletePayment(req, refPayment, "Successful cancel")
		completedPayment, err = s.storage.SavePayment(ctx, completedPayment, tx)
		if err != nil {
			customer.Balance = customer.Balance - req.Amount
			customer.BlockedMoney = customer.BlockedMoney + req.Amount
			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           customer.Id,
				Balance:      customer.Balance,
				BlockedMoney: customer.BlockedMoney,
			})

			merchant.BlockedMoney = merchant.BlockedMoney + req.Amount

			merchant, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      merchant.Id,
				Balance: merchant.BlockedMoney,
			})
			return nil, err
		}

		sts = append(sts, &authpb.StatementRequest{
			AccountId: customer.Id,
			PaymentId: completedPayment.PaymentId.String(),
		})

		sts = append(sts, &authpb.StatementRequest{
			AccountId: merchant.Id,
			PaymentId: completedPayment.PaymentId.String(),
		})
		// send statements to auth service
		if err := createStatement(ctx, s.client, sts); err != nil {
			customer.Balance = customer.Balance - req.Amount
			customer.BlockedMoney = customer.BlockedMoney + req.Amount
			customer, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:           customer.Id,
				Balance:      customer.Balance,
				BlockedMoney: customer.BlockedMoney,
			})

			merchant.BlockedMoney = merchant.BlockedMoney + req.Amount

			merchant, _ = s.client.UpdateBalance(ctx, &authpb.UpdateBalanceRequest{
				Id:      merchant.Id,
				Balance: merchant.BlockedMoney,
			})
			return nil, err
		}
		// commit tx
		if err := tx.Commit(); err != nil {
			return nil, err
		}
		return &paymentpb.Statement{
			PaymentId: completedPayment.PaymentId.String(),
			Status:    completedPayment.Status,
		}, nil
	}
	return &paymentpb.Statement{
		PaymentId: req.PaymentId,
		Status:    "Invalid transaction",
	}, nil
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
		_, err := stream.Recv()
		if err != nil {
			return err
		}
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	return nil
}
