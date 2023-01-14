package storage

import (
	"context"
	"database/sql"

	paymentpb "github.com/Edbeer/payment-grpc/proto"
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

func (s *PostgresStorage) SavePayment(ctx context.Context, payment *types.Payment) (*types.Payment, error) {
	query := `INSERT INTO payment (merchant, 
		customer, card_number, card_expiry_month,
		card_expiry_year, currency, operation,
		status, amount, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING *`
	pay := &types.Payment{}
	tx, err := s.db.BeginTx(ctx, nil)
	if tx != nil {
		defer tx.Rollback()
	}
	if err != nil {
		return nil, err
	}
	if err := s.db.QueryRowContext(
		ctx, query,
		payment.Merchant,
		payment.Customer,
		payment.CardNumber,
		payment.CardExpiryMonth,
		payment.CardExpiryYear,
		payment.Currency,
		payment.Operation,
		payment.Status,
		payment.Amount,
		payment.CreatedAt,
	).Scan(
		&pay.PaymentId, &pay.Merchant,
		&pay.Customer, &pay.CardNumber,
		&pay.CardExpiryMonth, &pay.CardExpiryYear, 
		&pay.Currency, &pay.Operation,
		&pay.Status, &pay.Amount,
		&pay.CreatedAt,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return pay, nil
}

func (s *PostgresStorage) GetPaymentByID(ctx context.Context, req *paymentpb.PaidRequest) (*types.Payment, error) {
	query := `SELECT * FROM payment WHERE payment_id = $1`
	pay := &types.Payment{}
	if err := s.db.QueryRowContext(
		ctx, query, req.PaymentId,
	).Scan(
		&pay.PaymentId, &pay.Merchant,
		&pay.Customer, &pay.CardNumber,
		&pay.CardExpiryMonth, &pay.CardExpiryYear, 
		&pay.Currency, &pay.Operation,
		&pay.Status, &pay.Amount,
		&pay.CreatedAt,
	); err != nil {
		return nil, err
	}
	return pay, nil
}