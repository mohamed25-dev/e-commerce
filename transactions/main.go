package main

import (
	"context"
	db "ecommerce/transactions/db/sqlc"
	rpc "ecommerce/transactions/proto"
	"ecommerce/transactions/server"
	"ecommerce/transactions/service"

	// TODO: rename to rpc and rename rpc to proto

	workflow "ecommerce/transactions/worfklow"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
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
	// DB connection
	//TODO: use env variables
	pgConn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/transactions_db")
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to DB >>>>>>>>>>")

	c, err := client.Dial(client.Options{})
	if err != nil {
		fmt.Println("Error while connecting to temporal")
	}
	defer c.Close()

	queries := db.New(pgConn)
	go startTemporalWorker(c, queries)

	transactionsService := &service.TransactionsService{Queries: *queries, TemporalClient: c}
	lis, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal("failed to listen, error: ", err)
	}

	streams := make(map[chan *rpc.CreateTransactionResponse]struct{})
	grpcServer := server.InitTransactionsRpcServer(transactionsService, &streams)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve, error: ", err)
	}

}
