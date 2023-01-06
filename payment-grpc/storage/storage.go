package storage

import (
	"context"
	"database/sql"

	"github.com/Edbeer/payment-grpc/types"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) SavePayment(ctx context.Context, tx *sql.Tx, payment *types.Payment) (*types.Payment, error) {
	query := `INSERT INTO payment (payment_id, payment_receiver, 
		payer, currency, operation, status, amount, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING *`
	pay := &types.Payment{}
	if err := tx.QueryRowContext(
		ctx, query,
		payment.PaymentId,
		payment.PaymentReceiver,
		payment.Payer,
		payment.Currency,
		payment.Operation,
		payment.Status,
		payment.Amount,
		payment.CreatedAt,
	).Scan(
		&pay.PaymentId, &pay.PaymentReceiver,
		&pay.Payer, &pay.Currency,
		&pay.Operation, &pay.Status,
		&pay.Amount, &pay.CreatedAt,
	); err != nil {
		return nil, err
	}
	return pay, nil
}


func (s *PostgresStorage) SaveBalance(ctx context.Context, tx *sql.Tx, account *types.Account, balance, bmoney uint64) (*types.Account, error) {
	query := `UPDATE account
				SET balance = COALESCE(NULLIF($1, 0), balance),
					blocked_money = COALESCE(NULLIF($2, 0), blocked_money)
				WHERE id = $3
				RETURNING *`
	acc := &types.Account{}
	if err := tx.QueryRowContext(
		ctx, query,
		balance,
		bmoney,
		account.ID,
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, &acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	return acc, nil
}
