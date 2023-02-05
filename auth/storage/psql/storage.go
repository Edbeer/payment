package storage

import (
	"context"
	"database/sql"

	authpb "github.com/Edbeer/auth/proto"
	"github.com/Edbeer/auth/types"
	"github.com/lib/pq"
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

// Create account
func (s *PostgresStorage) CreateAccount(ctx context.Context, account *authpb.CreateRequest) (*types.Account, error) {
	query := `INSERT INTO account (first_name, 
		last_name, card_number, card_expiry_month, 
		card_expiry_year, card_security_code, 
		balance, blocked_money, statement, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
			RETURNING *`
	reqAcc := types.NewAccount(account)
	acc := &types.Account{}
	if err := s.db.QueryRowContext(
		ctx, query,
		reqAcc.FirstName,
		reqAcc.LastName,
		reqAcc.CardNumber,
		reqAcc.CardExpiryMonth,
		reqAcc.CardExpiryYear,
		reqAcc.CardSecurityCode,
		reqAcc.Balance,
		reqAcc.BlockedMoney,
		pq.Array(reqAcc.Statement),
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, pq.Array(&acc.Statement), 
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}

	return acc, nil
}

// Update Account
func (s *PostgresStorage) UpdateAccount(ctx context.Context, account *authpb.UpdateRequest) (*types.Account, error) {
	query := `UPDATE account
				SET first_name = COALESCE(NULLIF($1, ''), first_name),
					last_name = COALESCE(NULLIF($2, ''), last_name),
					card_number = COALESCE(NULLIF($3, ''), card_number),
					card_expiry_month = COALESCE(NULLIF($4, ''), card_expiry_month),
					card_expiry_year = COALESCE(NULLIF($5, ''), card_expiry_year),
					card_security_code = COALESCE(NULLIF($6, ''), card_security_code)
				WHERE id = $7
				RETURNING *`
	acc := &types.Account{}
	if err := s.db.QueryRowContext(
		ctx, query,
		account.FirstName,
		account.LastName,
		account.CardNumber,
		account.CardExpiryMonth,
		account.CardExpiryYear,
		account.CardSecurityCode,
		account.Id,
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, pq.Array(&acc.Statement), 
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	return acc, nil
}

// Delete Account
func (s *PostgresStorage) DeleteAccount(ctx context.Context, req *authpb.DeleteRequest) (*authpb.DeleteResponse, error) {
	query := `DELETE FROM account WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, req.Id)
	if err != nil {
		return nil, err
	}
	return &authpb.DeleteResponse{
		Status: "Account was deleted",
	}, err
}

// Deposit account
func (s *PostgresStorage) DepositAccount(ctx context.Context, req *authpb.DepositRequest) (*authpb.DepositResponse, error) {
	query := `UPDATE account
		SET balance = COALESCE(NULLIF($1, 0), balance)
		WHERE card_number = $2
		RETURNING *`
	acc := &types.Account{}

	if err := s.db.QueryRowContext(
		ctx,
		query,
		req.Balance,
		req.CardNumber,
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, pq.Array(&acc.Statement), 
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &authpb.DepositResponse{
		Status: "Successful deposit",
	}, nil
}

func (s *PostgresStorage) GetAccount(ctx context.Context) ([]*types.Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		acc := &types.Account{}
		if err := rows.Scan(
			&acc.ID, &acc.FirstName,
			&acc.LastName, &acc.CardNumber,
			&acc.CardExpiryMonth, &acc.CardExpiryYear,
			&acc.CardSecurityCode, &acc.Balance,
			&acc.BlockedMoney, pq.Array(&acc.Statement), 
			&acc.CreatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStorage) GetAccountByID(ctx context.Context, req *authpb.GetIDRequest) (*types.Account, error) {
	query := `SELECT * FROM account 
			WHERE id = $1`
	acc := &types.Account{}

	if err := s.db.QueryRowContext(
		ctx, query, req.Id,
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, pq.Array(&acc.Statement), 
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *PostgresStorage) SaveBalance(ctx context.Context, req *authpb.UpdateBalanceRequest) (*types.Account, error) {
	query := `UPDATE account
				SET balance = COALESCE($1, balance),
					blocked_money = COALESCE($2, blocked_money)
				WHERE id = $3
				RETURNING *`
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	acc := &types.Account{}
	if err := tx.QueryRowContext(
		ctx, query,
		req.Balance,
		req.BlockedMoney,
		req.Id,
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, pq.Array(&acc.Statement), 
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return  nil, err
	}
	return acc, nil
}

func (s *PostgresStorage) UpdateStatement(ctx context.Context, req *authpb.StatementRequest) ([]string, error) {
	query := `UPDATE account
				SET statement = array_append(statement, $1)
				WHERE id = $2
				RETURNING *`
	acc := &types.Account{}
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	if err := tx.QueryRowContext(
		ctx,
		query,
		req.PaymentId,
		req.AccountId,
	).Scan(
		&acc.ID, &acc.FirstName,
		&acc.LastName, &acc.CardNumber,
		&acc.CardExpiryMonth, &acc.CardExpiryYear,
		&acc.CardSecurityCode, &acc.Balance,
		&acc.BlockedMoney, pq.Array(&acc.Statement),
		&acc.CreatedAt,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return acc.Statement, nil
}