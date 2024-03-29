package workflow

import (
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type TransactionWorkflow struct {
	Activities *TransactionActivity
}

func (w *TransactionWorkflow) CreateTransactionWorkflow(ctx workflow.Context, data models.CreateTransactionRequestModel, product db.Product, customer db.Customer) (db.Transaction, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 1,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var createdTransaction db.Transaction
	err := workflow.ExecuteActivity(ctx, w.Activities.CreateTransactionActivity, data).Get(ctx, &createdTransaction)
	if err != nil {
		return db.Transaction{}, err
	}

	analyticsTransactionData, err := models.MapTransactionDataToAnalyticsTransaction(data, product, customer)
	if err != nil {
		fmt.Println("mapping to analytics transaction failed, err: ", err)
		return db.Transaction{}, err
	}

	err = workflow.ExecuteActivity(ctx, w.Activities.SendTransactionToAnalyticsActivity, analyticsTransactionData).Get(ctx, nil)
	if err != nil {
		return db.Transaction{}, err
	}

	return createdTransaction, err
}
