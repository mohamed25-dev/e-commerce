package db

import (
	"context"
	"ecommerce/transactions/utils"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func generateRandomProduct() CreateProductParams {
	val, err := utils.Float32ToPgNumeric(float32(utils.GenerateRandomAmount()))
	if err != nil {
		log.Fatal("converting to pg numeric failed")
	}

	return CreateProductParams{
		ID:          uuid.NewString(),
		ProductName: utils.GenerateRandomString(8),
		Price:       val,
	}
}

func TestGetProductById(t *testing.T) {
	generatedProduct := generateRandomProduct()
	createdProduct, err := testQueries.CreateProduct(context.Background(), generatedProduct)
	if err != nil {
		t.Fatal(err)
	}

	product, err := testQueries.GetProductById(context.Background(), createdProduct.ID)

	require.NoError(t, err)
	require.Equal(t, createdProduct.ID, product.ID)
	require.Equal(t, createdProduct.ProductName, product.ProductName)
	require.Equal(t, createdProduct.Price, product.Price)
}
