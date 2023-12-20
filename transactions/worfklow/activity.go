package workflow

import (
	"context"
	analyticsRpc "ecommerce/analytics/proto"
	rpc "ecommerce/transactions/clients"
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"fmt"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type TransactionActivity struct {
	Queries *db.Queries
}

func (a *TransactionActivity) CreateTransactionActivity(ctx context.Context, data models.CreateTransactionRequestModel) (db.Transaction, error) {
	transactionParams, err := models.MapTransactionDataToDbParams(data)
	transactionParams.ID = uuid.NewString()

	if err != nil {
		fmt.Println("mapping to transaction params failed, err: ", err)
		return db.Transaction{}, err
	}

	createdTransaction, err := a.Queries.CreateTransaction(ctx, transactionParams)
	if err != nil {
		log.Println("error while creating transaction: ", err)
		return db.Transaction{}, err
	}

	return createdTransaction, nil
}

func (a *TransactionActivity) SendTransactionToAnalyticsActivity(ctx context.Context, createdTransaction models.CreateAnalyticsTransactionRequestModel) error {
	//TODO: use env variables
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Println("connection failed: ", err)
	}
	defer conn.Close()

	// TODO: check which is better, creating a connection each time, or keeping the connection open somewhere
	// NOTE: In my opinion using async way of communication here might be a better choice, for example using a message
	// queuing system like RapidMQ to decouple the services from each other. However I wen with the sync communication as
	// this is a good use case to test temporal workflows. Additionally it guarntees consistency and realtime analytics
	c := analyticsRpc.NewAnalticsServiceClient(conn)
	err = rpc.CreateAnalyticsTransaction(ctx, c, createdTransaction)
	if err != nil {
		log.Println("something went wrong while creating a transaction in the analytics service, err: ", err)
		return err
	}

	return nil
}
