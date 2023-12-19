package service

import (
	"context"
	db "ecommerce/analytics/db/sqlc"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kataras/go-events"
)

type AnalyticsService struct {
	Queries               db.Queries
	EventEmmiter          events.EventEmmiter
	Mu                    sync.Mutex
	TopCustomerStreams    map[chan *[]db.GetTopCustomersRow]struct{}
	TotalSalesStreams     map[chan *db.GetTotalSalesRow]struct{}
	SalesByProductStreams map[chan *db.GetTotalSalesByProductIdRow]struct{}
}

func (service *AnalyticsService) GetTransactionById(ctx context.Context, id string) (db.Transaction, error) {
	trx, err := service.Queries.GetTransactionById(ctx, id)
	if err != nil {
		// check the error type, maybe return transaction not found
		return db.Transaction{}, err
	}

	return trx, err
}

func (service *AnalyticsService) CreateTransaction(ctx context.Context, transactionData db.CreateTransactionParams) (db.Transaction, error) {
	transactionData.ID = uuid.NewString()
	createdTransaction, err := service.Queries.CreateTransaction(ctx, transactionData)
	if err != nil {
		return db.Transaction{}, err
	}

	service.EventEmmiter.Emit("transaction_created", createdTransaction)
	return createdTransaction, err
}

// NOTE: performance of this function can be enhanced by using multiple goroutines to
// handle the streaming simoltanously
func (service *AnalyticsService) HandleTransactionCreatedEvent(data ...interface{}) {
	var createdTransaction db.Transaction
	var ok bool

	if createdTransaction, ok = data[0].(db.Transaction); !ok {
		fmt.Println("converting received data to transaction failed.")
		return
	}

	totalSales, err := service.Queries.GetTotalSales(context.Background())
	if err != nil {
		fmt.Println("something went wrong while retrieving total sales, err: ", err)
		return
	}

	salesByProduct, err := service.Queries.GetTotalSalesByProductId(context.Background(), createdTransaction.ProductID)
	if err != nil {
		fmt.Println("something went wrong while retrieving total sales, err: ", err)
		return
	}

	topCustomers, err := service.Queries.GetTopCustomers(context.Background(), 5)
	if err != nil {
		fmt.Println("something went wrong while retrieving top customers, err: ", err)
		return
	}

	service.Mu.Lock()
	for ch := range service.TopCustomerStreams {
		ch <- &topCustomers
	}

	for ch := range service.TotalSalesStreams {
		ch <- &totalSales
	}

	for ch := range service.SalesByProductStreams {
		ch <- &salesByProduct
	}
	service.Mu.Unlock()
}
