package db

import (
	"context"
	"ecommerce/transactions/utils"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func generateRandomTransaction() CreateTransactionParams {
	val, err := utils.Float32ToPgNumeric(float32(utils.GenerateRandomAmount()))
	if err != nil {
		log.Fatal("converting to pg numeric failed")
	}

	transaction := CreateTransactionParams{
		ID:         uuid.NewString(),
		ProductID:  uuid.NewString(),
		CustomerID: uuid.NewString(),
		Quantity:   int32(utils.GenerateRandomInt(1, 20)),
		TotalPrice: val,
	}

	return transaction
}

func TestGetTransactionById(t *testing.T) {
	generatedTransaction := generateRandomTransaction()
	createdTransaction, err := testQueries.CreateTransaction(context.Background(), generatedTransaction)
	if err != nil {
		t.Fatal(err)
	}

	transaction, err := testQueries.GetTransactionById(context.Background(), createdTransaction.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, createdTransaction.ID, transaction.ID)
	require.Equal(t, createdTransaction.CustomerID, transaction.CustomerID)
	require.Equal(t, createdTransaction.ProductID, transaction.ProductID)
	require.Equal(t, createdTransaction.Quantity, transaction.Quantity)
	require.Equal(t, createdTransaction.TotalPrice, transaction.TotalPrice)
	require.NotZero(t, createdTransaction.CreatedAt)
}

func TestCreateTransaction(t *testing.T) {
	val, err := utils.Float32ToPgNumeric(float32(utils.GenerateRandomAmount()))
	if err != nil {
		t.Fatal("converting to pg numeric failed")
	}

	transaction := CreateTransactionParams{
		ID:         uuid.NewString(),
		ProductID:  uuid.NewString(),
		CustomerID: uuid.NewString(),
		Quantity:   int32(utils.GenerateRandomInt(1, 20)),
		TotalPrice: val,
	}

	createdTransaction, err := testQueries.CreateTransaction(context.Background(), transaction)

	require.NoError(t, err)
	require.NotEmpty(t, transaction)
	require.Equal(t, createdTransaction.ID, transaction.ID)
	require.Equal(t, createdTransaction.CustomerID, transaction.CustomerID)
	require.Equal(t, createdTransaction.ProductID, transaction.ProductID)
	require.Equal(t, createdTransaction.Quantity, transaction.Quantity)
	require.Equal(t, createdTransaction.TotalPrice, transaction.TotalPrice)
	require.NotZero(t, createdTransaction.CreatedAt)
}
