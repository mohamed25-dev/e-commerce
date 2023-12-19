package main

import (
	"context"
	db "ecommerce/analytics/db/sqlc"
	"ecommerce/analytics/server"
	"ecommerce/analytics/service"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/kataras/go-events"
)

func main() {
	//TODO: use env variables
	pgConn, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/analytics_db")
	if err != nil {
		panic(err)
	}

	events := events.New()
	queries := db.New(pgConn)

	topCustomerStreams := make(map[chan string]struct{})
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

	//TODO: use env variables
	lis, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatal("failed to listen, error: ", err)
	}

	grpcServer := server.InitAnalyticsRpcServer(analyticsService, totalSalesStreams, topCustomerStreams, salesByProductStreams)

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal("failed to serve, error: ", err)
	}
}
