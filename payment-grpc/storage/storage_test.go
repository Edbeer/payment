package storage

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	authpb "github.com/Edbeer/auth-grpc/proto"
	paymentpb "github.com/Edbeer/payment-grpc/proto"
	"github.com/Edbeer/payment-grpc/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_SavePayment(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("SavePayment", func(t *testing.T) {
		customer := &authpb.Account{
			Id:               uuid.New().String(),
			FirstName:        "Pasha",
			LastName:         "Volkov",
			CardNumber:       "444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "924",
			Balance:          50,
			BlockedMoney:     0,
			CreatedAt:       timestamppb.Now(),
		}

		merchant := &authpb.Account{
			Id:               uuid.New().String(),
			FirstName:        "Pasha1",
			LastName:         "Volkov1",
			CardNumber:       "444444444444443",
			CardExpiryMonth:  "10",
			CardExpiryYear:   "25",
			CardSecurityCode: "923",
			Balance:          0,
			BlockedMoney:     0,
			CreatedAt:       timestamppb.Now(),
		}

		payReq := &paymentpb.CreateRequest{
			Merchant:         merchant.Id,
			Customer:         customer.Id,
			CardNumber:       "444444444444444",
			CardExpiryMonth:  "12",
			CardExpiryYear:   "24",
			CardSecurityCode: "924",
			Currency:         "RUB",
			Amount:           50,
		}
		
		payment := types.CreateAuthPayment(payReq, customer, merchant, "")

		colums := []string{
			"payment_id",
			"merchant",
			"customer",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"currency",
			"operation",
			"status",
			"amount",
			"created_at",
		}
		rows := sqlmock.NewRows(colums).AddRow(
			payment.PaymentId,
			merchant.Id,
			customer.Id,
			"444444444444444",
			"12",
			"24",
			"RUB",
			payment.Operation,
			"",
			50,
			payment.CreatedAt,
		)
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO payment (merchant, 
			customer, card_number, card_expiry_month,
			card_expiry_year, currency, operation,
			status, amount, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING *`)).WithArgs(
					payment.Merchant,
					payment.Customer,
					payment.CardNumber,
					payment.CardExpiryMonth,
					payment.CardExpiryYear,
					payment.Currency,
					payment.Operation,
					payment.Status,
					payment.Amount,
					payment.CreatedAt,).WillReturnRows(rows)
		tx, _ := db.BeginTx(context.Background(), nil)
		pay, err := psql.SavePayment(context.Background(), payment, tx)
		require.NoError(t, err)
		require.NotNil(t, pay)
	})
}

func Test_GetPaymentByID(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	psql := NewPostgresStorage(db)

	t.Run("GetPaymentByID", func(t *testing.T) {
		req := &paymentpb.PaidRequest{
			PaymentId: uuid.New().String(),
		}

		colums := []string{
			"payment_id",
			"merchant",
			"customer",
			"card_number",
			"card_expiry_month",
			"card_expiry_year",
			"currency",
			"operation",
			"status",
			"amount",
			"created_at",
		}
		rows := sqlmock.NewRows(colums).AddRow(
			req.PaymentId,
			"",
			"",
			"444444444444444",
			"12",
			"24",
			"RUB",
			"",
			"",
			50,
			time.Now(),
		)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM payment WHERE payment_id = $1`)).WithArgs(req.PaymentId).WillReturnRows(rows)

		pay, err := psql.GetPaymentByID(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, pay)
	})
}