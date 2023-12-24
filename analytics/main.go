package main

import (
	"context"
	server "ecommerce/analytics/api"
	db "ecommerce/analytics/db/sqlc"
	"ecommerce/analytics/service"
	"ecommerce/transactions/utils"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/kataras/go-events"
)

func main() {
	err := godotenv.Load()
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

	events := events.New()
	queries := db.New(pool)

	topCustomerStreams := make(map[chan *[]db.GetTopCustomersRow]struct{})
	totalSalesStreams := make(map[chan *db.GetTotalSalesRow]struct{})
	salesByProductStreams := make(map[chan *db.GetTotalSalesByProductIdRow]struct{})

	analyticsService := &service.AnalyticsService{
		Queries:               *queries,
		EventEmmiter:          events,
		TopCustomerStreams:    topCustomerStreams,
		TotalSalesStreams:     totalSalesStreams,
		SalesByProductStreams: salesByProductStreams,
	}

	events.On("transaction_created", analyticsService.HandleTransactionCreatedEvent)

	lis, err := net.Listen("tcp", os.Getenv("ANALYTICS_SERVICE_IP"))
	if err != nil {
		log.Fatal("failed to listen, error: ", err)
	}

	grpcServer := server.InitAnalyticsRpcServer(analyticsService, totalSalesStreams, topCustomerStreams, salesByProductStreams)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve, error: ", err)
	}
}
