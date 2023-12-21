package main

import (
	"context"
	server "ecommerce/transactions/api"
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/proto"
	"ecommerce/transactions/service"
	"ecommerce/transactions/utils"
	"os"

	workflow "ecommerce/transactions/workflow"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func startTemporalWorker(c client.Client, q *db.Queries) error {
	w := worker.New(c, "transactions-queue", worker.Options{})
	ta := workflow.TransactionActivity{Queries: q}
	tw := workflow.TransactionWorkflow{Activities: &ta}

	w.RegisterWorkflow(tw.CreateTransactionWorkflow)
	w.RegisterActivity(ta.SendTransactionToAnalyticsActivity)
	w.RegisterActivity(ta.CreateTransactionActivity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Println("unable to start worker, er: ", err)
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load("../transactions/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr, err := utils.GetDbConnectionString()
	if err != nil {
		log.Fatal("could not get DB connection string, err: ", err)
	}

	// Create a connection pool
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close(context.Background())

	queries := db.New(pool)

	c, err := client.Dial(client.Options{})
	if err != nil {
		fmt.Println("Error while connecting to temporal")
	}
	defer c.Close()
	go startTemporalWorker(c, queries)

	transactionsService := &service.TransactionsService{Queries: queries, TemporalClient: c}
	lis, err := net.Listen("tcp", os.Getenv("TRANSACTIONS_SERVICE_IP"))
	if err != nil {
		log.Fatal("failed to listen, error: ", err)
	}

	streams := make(map[chan *proto.CreateTransactionResponse]struct{})
	grpcServer := server.InitTransactionsRpcServer(transactionsService, &streams)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve, error: ", err)
	}

}
