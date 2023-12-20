package clients

import (
	"context"
	analyticsrpc "ecommerce/analytics/proto"
	"ecommerce/transactions/models"
	"log"
)

type CreateTransactionRequest struct {
	CustomerId  string
	ProductId   string
	Quantity    int32
	TotalAmount float32
}

func CreateAnalyticsTransaction(ctx context.Context, c analyticsrpc.AnalticsServiceClient, createdTransaction models.CreateAnalyticsTransactionRequestModel) error {
	requestBody := &analyticsrpc.CreateAnalyticsTransactionRequest{
		CustomerId:   createdTransaction.CustomerID,
		CustomerName: createdTransaction.CustomerName,
		ProductId:    createdTransaction.ProductID,
		ProductName:  createdTransaction.ProductName,
		Quantity:     createdTransaction.Quantity,
		TotalAmount:  createdTransaction.TotalPrice,
	}

	_, err := c.CreateAnalyticsTransaction(ctx, requestBody)
	if err != nil {
		log.Println("did not get response, err: ", err)
		return err
	}

	return nil
}
