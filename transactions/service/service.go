package service

import (
	"context"
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"ecommerce/transactions/utils"
	workflow "ecommerce/transactions/workflow"
	"errors"
	"log"

	"go.temporal.io/sdk/client"
)

type TransactionServicer interface {
	GetTransactionById(ctx context.Context, id string) (db.Transaction, error)
	CreateTransaction(ctx context.Context, transactionData models.CreateTransactionRequestModel) (db.Transaction, error)
}

type TransactionsService struct {
	Queries        db.Querier
	TemporalClient client.Client
}

func (service *TransactionsService) GetTransactionById(ctx context.Context, id string) (db.Transaction, error) {
	trx, err := service.Queries.GetTransactionById(ctx, id)
	if err != nil {
		return db.Transaction{}, errors.New("transaction not found")
	}

	return trx, err
}

func (service *TransactionsService) CreateTransaction(ctx context.Context, transactionData models.CreateTransactionRequestModel) (db.Transaction, error) {
	customer, err := service.Queries.GetCustomerById(ctx, transactionData.CustomerID)
	if err != nil {
		return db.Transaction{}, errors.New("customer not found")
	}

	product, err := service.Queries.GetProductById(ctx, transactionData.ProductID)
	if err != nil {
		return db.Transaction{}, errors.New("product not found")
	}

	productPrice, err := utils.PgNumericToFloat32(product.Price)
	if err != nil {
		return db.Transaction{}, err
	}

	transactionData.TotalPrice = float32(transactionData.Quantity) * productPrice

	options := client.StartWorkflowOptions{
		ID:        "transactions-queue",
		TaskQueue: "transactions-queue", // Task queue for the workflow
	}

	activities := workflow.TransactionActivity{
		Queries: service.Queries,
	}

	workflow := workflow.TransactionWorkflow{Activities: &activities}
	run, err := service.TemporalClient.ExecuteWorkflow(ctx, options, workflow.CreateTransactionWorkflow, transactionData, product, customer)
	if err != nil {
		log.Println("error while executing temporal workflow: ", err)
		return db.Transaction{}, err
	}

	var createdTransaction db.Transaction
	err = run.Get(ctx, &createdTransaction)
	if err != nil {
		log.Println("error while get results from temporal workflow: ", err)
		return db.Transaction{}, err
	}

	return createdTransaction, err
}
