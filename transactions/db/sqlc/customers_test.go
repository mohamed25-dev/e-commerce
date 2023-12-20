package db

import (
	"context"
	"ecommerce/transactions/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func generateRandomCustomer() CreateCustomerParams {
	return CreateCustomerParams{
		ID:           uuid.NewString(),
		CustomerName: utils.GenerateRandomString(8),
	}
}

func TestGetCustomerById(t *testing.T) {
	generatedCustomer := generateRandomCustomer()
	createdCustomer, err := testQueries.CreateCustomer(context.Background(), generatedCustomer)
	if err != nil {
		t.Fatal(err)
	}

	customer, err := testQueries.GetCustomerById(context.Background(), createdCustomer.ID)

	require.NoError(t, err)
	require.Equal(t, createdCustomer.ID, customer.ID)
	require.Equal(t, createdCustomer.CustomerName, customer.CustomerName)
}
