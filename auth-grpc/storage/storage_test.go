package storage

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	authpb "github.com/Edbeer/auth-grpc/proto"
	"github.com/Edbeer/auth-grpc/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_CreateAccount(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("Create", func(t *testing.T) {
		req := &authpb.CreateRequest{
			FirstName:        "Pasha",
			LastName:         "Volkov",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
		}

		account := types.NewAccount(req)

		colums := []string{
			"id",
			"first_name",
			"last_name",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"card_security_code",
			"balance", "blocked_money",
			"created_at",
		}
		rows := sqlmock.NewRows(colums).AddRow(
			account.ID,
			"Pasha",
			"Volkov",
			"4444444444444444",
			"12",
			"24",
			"123",
			0,
			0,
			account.CreatedAt,
		)
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO account (first_name, 
			last_name, card_number, card_expiry_month, 
			card_expiry_year, card_security_code, 
			balance, blocked_money, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
				RETURNING *`)).WithArgs(
			account.FirstName,
			account.LastName,
			account.CardNumber,
			account.CardExpiryMonth,
			account.CardExpiryYear,
			account.CardSecurityCode,
			account.Balance,
			account.BlockedMoney,).WillReturnRows(rows)
		createdUser, err := psql.CreateAccount(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, createdUser)
		require.Equal(t, createdUser, account)	
	})
}

func Test_UpdateAccount(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("Update", func(t *testing.T) {
		reqToCreate := &authpb.CreateRequest{
			FirstName:        "Pasha",
			LastName:         "Volkov",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
		}

		account := types.NewAccount(reqToCreate)

		reqToUpdate := &authpb.UpdateRequest{
			FirstName:        "Pasha1",
			LastName:         "",
			CardNumber:       "4444444444444443",
			CardExpiryMonth:  "",
			CardExpiryYear:   "",
			CardSecurityCode: "",
			Id: account.ID.String(),
		}
		colums := []string{
			"id",
			"first_name",
			"last_name",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"card_security_code",
			"balance", "blocked_money",
			"created_at",
		}
		rows := sqlmock.NewRows(colums).AddRow(
			account.ID,
			"Pasha1",
			"Volkov",
			"4444444444444443",
			"12",
			"24",
			"123",
			0,
			0,
			account.CreatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE account
		SET first_name = COALESCE(NULLIF($1, ''), first_name),
			last_name = COALESCE(NULLIF($2, ''), last_name),
			card_number = COALESCE(NULLIF($3, ''), card_number),
			card_expiry_month = COALESCE(NULLIF($4, ''), card_expiry_month),
			card_expiry_year = COALESCE(NULLIF($5, ''), card_expiry_year),
			card_security_code = COALESCE(NULLIF($6, ''), card_security_code)
		WHERE id = $7
		RETURNING *`)).WithArgs(
			reqToUpdate.FirstName,
			reqToUpdate.LastName,
			reqToUpdate.CardNumber,
			reqToUpdate.CardExpiryMonth,
			reqToUpdate.CardExpiryYear,
			reqToUpdate.CardSecurityCode,
			reqToUpdate.Id).WillReturnRows(rows)

		updatedAccount, err := psql.UpdateAccount(context.Background(), reqToUpdate)
		require.NoError(t, err)
		require.NotEqual(t, updatedAccount, account)
	})
}

func Test_DeleteAccount(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("Delete", func(t *testing.T) {
		uid := uuid.New().String()
		req := &authpb.DeleteRequest{
			Id: uid,
		}
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM account WHERE id = $1`)).WithArgs(uid).WillReturnResult(sqlmock.NewResult(1, 1))

		_, err := psql.DeleteAccount(context.Background(), req)
		require.NoError(t, err)
	})
}

func Test_DepositAccount(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("Deposit", func(t *testing.T) {
		req := &authpb.DepositRequest{
			CardNumber:       "4444444444444444",
			Balance: 50,
		}

		reqToCreate := &authpb.CreateRequest{
			FirstName:        "Pasha",
			LastName:         "Volkov",
			CardNumber:       "4444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "123",
		}

		account := types.NewAccount(reqToCreate)

		colums := []string{
			"id",
			"first_name",
			"last_name",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"card_security_code",
			"balance", "blocked_money",
			"created_at",
		}
		rows := sqlmock.NewRows(colums).AddRow(
			account.ID,
			"Pasha",
			"Volkov",
			"4444444444444444",
			"12",
			"24",
			"123",
			50,
			0,
			account.CreatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE account
		SET balance = COALESCE(NULLIF($1, 0), balance)
		WHERE card_number = $2
		RETURNING *`)).WithArgs(uint64(50), req.CardNumber).WillReturnRows(rows)
		_, err := psql.DepositAccount(context.Background(), req)
		require.NoError(t, err)
		require.Equal(t, req.CardNumber, account.CardNumber)
		require.Equal(t, req.Balance, uint64(50))
	})
}

func Test_GetAccount(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("GetAccounts", func(t *testing.T) {
		req1 := &authpb.CreateRequest{
			FirstName:        "Pasha1",
			LastName:         "volkov1",
			CardNumber:       "444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "924",
		}
		account1 := types.NewAccount(req1)
		colums := []string{
			"id",
			"first_name",
			"last_name",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"card_security_code",
			"balance", "blocked_money",
			"created_at",
		}
		rows1 := sqlmock.NewRows(colums).AddRow(
			account1.ID,
			"Pasha1",
			"volkov1",
			"444444444444444",
			"12",
			"24",
			"924",
			0,
			0,
			account1.CreatedAt,
		)
		req2 := &authpb.CreateRequest{
			FirstName:        "Pasha",
			LastName:         "volkov",
			CardNumber:       "444444444444432",
			CardExpiryMonth:  "10",
			CardExpiryYear:   "25",
			CardSecurityCode: "934",
		}
		account2 := types.NewAccount(req2)

		rows2 := sqlmock.NewRows(colums).AddRow(
			account2.ID,
			"Pasha1",
			"volkov1",
			"444444444444444",
			"12",
			"24",
			"924",
			0,
			0,
			account2.CreatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM account`)).WillReturnRows(rows1, rows2)
		userList, err := psql.GetAccount(context.Background())
		require.NoError(t, err)
		require.NotNil(t, userList)
	})
}

func Test_GetAccountByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("GetAccountByID", func(t *testing.T) {
		req := &authpb.CreateRequest{
			FirstName:        "Pasha1",
			LastName:         "volkov1",
			CardNumber:       "444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "924",
		}
		account := types.NewAccount(req)
		colums := []string{
			"id",
			"first_name",
			"last_name",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"card_security_code",
			"balance", "blocked_money",
			"created_at",
		}
		rows := sqlmock.NewRows(colums).AddRow(
			account.ID,
			"Pasha1",
			"volkov1",
			"444444444444444",
			"12",
			"24",
			"924",
			0,
			0,
			account.CreatedAt,
		)
		reqID := &authpb.GetIDRequest{
			Id: account.ID.String(),
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM account WHERE id = $1`)).WithArgs(reqID.Id).WillReturnRows(rows)
		acc, err := psql.GetAccountByID(context.Background(), reqID)
		require.NoError(t, err)
		require.NotNil(t, acc)
	})
}