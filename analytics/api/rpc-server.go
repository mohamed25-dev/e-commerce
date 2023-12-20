package server

import (
	"context"
	db "ecommerce/analytics/db/sqlc"
	rpc "ecommerce/analytics/proto"
	"ecommerce/analytics/service"
	"ecommerce/transactions/utils"
	"fmt"
	"log"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

type AnalyticsServer struct {
	rpc.UnimplementedAnalticsServiceServer
	service *service.AnalyticsService
	// NOTE: creating two different mutexes might be a good idea,
	// as it will reduce the blocked code area for each mutex
	mu                    sync.Mutex
	topCustomerStreams    map[chan *[]db.GetTopCustomersRow]struct{}
	totalSalesStreams     map[chan *db.GetTotalSalesRow]struct{}
	salesByProductStreams map[chan *db.GetTotalSalesByProductIdRow]struct{}
}

func InitAnalyticsRpcServer(service *service.AnalyticsService, totalSalesStreams map[chan *db.GetTotalSalesRow]struct{}, topCustomerStreams map[chan *[]db.GetTopCustomersRow]struct{}, salesByProductStreams map[chan *db.GetTotalSalesByProductIdRow]struct{}) *grpc.Server {
	grpcServer := grpc.NewServer()

	analyticServer := &AnalyticsServer{
		service:               service,
		topCustomerStreams:    topCustomerStreams,
		totalSalesStreams:     totalSalesStreams,
		salesByProductStreams: salesByProductStreams,
	}

	rpc.RegisterAnalticsServiceServer(grpcServer, analyticServer)
	return grpcServer
}

func (s *AnalyticsServer) CreateAnalyticsTransaction(ctx context.Context, req *rpc.CreateAnalyticsTransactionRequest) (*rpc.CreateAnaltyicsTransactionResponse, error) {
	totalPrice, err := utils.Float32ToPgNumeric(req.TotalAmount)
	if err != nil {
		return nil, err
	}

	transactionData := db.CreateTransactionParams{
		CustomerID:   req.CustomerId,
		CustomerName: req.CustomerName,
		ProductID:    req.ProductId,
		ProductName:  req.ProductName,
		Quantity:     req.Quantity,
		TotalPrice:   totalPrice,
	}

	createdTransaction, err := s.service.CreateTransaction(ctx, transactionData)
	if err != nil {
		log.Println("something went wrong while creating the transaction, err: ", err)
		return nil, err
	}

	amountString := createdTransaction.TotalPrice.Int.String()
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		log.Println("Something went wrong while converting string to number, ", err)
		return nil, err
	}

	response := &rpc.CreateAnaltyicsTransactionResponse{
		Transaction: &rpc.AnalyticsTransaction{
			Id:          createdTransaction.ID,
			CustomerId:  createdTransaction.CustomerID,
			ProductId:   createdTransaction.ProductID,
			Quantity:    createdTransaction.Quantity,
			TotalAmount: float32(amount),
		},
	}

	return response, nil
}

func (s *AnalyticsServer) StreamTotalSales(empty *rpc.EmptyRequest, stream rpc.AnalticsService_StreamTotalSalesServer) error {
	totalSales, err := s.service.Queries.GetTotalSales(context.Background())
	if err != nil {
		fmt.Println("something went wrong wile retrieving data, err: ", err)
		return err
	}

	err = streamTotalSales(&totalSales, stream)
	if err != nil {
		fmt.Println("something went wrong wile streaming data, err: ", err)
		return err
	}

	s.mu.Lock()
	ch := make(chan *db.GetTotalSalesRow)
	s.totalSalesStreams[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.totalSalesStreams, ch)
		close(ch)
		s.mu.Unlock()
	}()

	for {
		totalSales, ok := <-ch
		if !ok {
			fmt.Println("channel is no longer available !!")
			break
		}

		err := streamTotalSales(totalSales, stream)
		if err != nil {
			fmt.Println("something went wrong while streaming data, err: ", err)
			return err
		}

	}

	return nil
}

func (s *AnalyticsServer) StreamSalesByProduct(req *rpc.StreamSalesByProductRequest, stream rpc.AnalticsService_StreamSalesByProductServer) error {
	salesByProduct, err := s.service.Queries.GetTotalSalesByProductId(context.Background(), req.ProductId)
	if err != nil {
		fmt.Println("fetching sales by product from DB failed, err: ", err)
	}

	err = streamSalesByProduct(req.ProductId, &salesByProduct, stream)
	if err != nil {
		return err
	}

	s.mu.Lock()
	ch := make(chan *db.GetTotalSalesByProductIdRow)
	s.salesByProductStreams[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.salesByProductStreams, ch)
		close(ch)
		s.mu.Unlock()
	}()

	for {
		salesByProduct, ok := <-ch
		if !ok {
			fmt.Println("channel is no longer available !!")
			break
		}

		streamSalesByProduct(req.ProductId, salesByProduct, stream)
	}

	return nil
}

func (s *AnalyticsServer) StreamTopCustomers(req *rpc.EmptyRequest, stream rpc.AnalticsService_StreamTopCustomersServer) error {
	fetchedTopCustomers, err := s.service.Queries.GetTopCustomers(context.Background(), 5)
	if err != nil {
		fmt.Println("fetching top customer from DB failed, err: ", err)
	}

	err = streamTopCustomers(&fetchedTopCustomers, stream)
	if err != nil {
		fmt.Println("error while streaming data, err: ", err)
	}

	s.mu.Lock()
	ch := make(chan *[]db.GetTopCustomersRow)
	s.topCustomerStreams[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.topCustomerStreams, ch)
		close(ch)
		s.mu.Unlock()
	}()

	for {
		receivedTopCustomers, ok := <-ch
		if !ok {
			fmt.Println("channel is no longer available !!")
			break
		}

		err := streamTopCustomers(receivedTopCustomers, stream)
		if err != nil {
			fmt.Println("error while streaming data, err: ", err)
		}
	}

	return nil
}

func streamTopCustomers(receivedTopCustomers *[]db.GetTopCustomersRow, stream rpc.AnalticsService_StreamTopCustomersServer) error {
	var customers []*rpc.TopCustomer
	for _, customer := range *receivedTopCustomers {
		totalPrice, err := utils.PgNumericToFloat32(customer.TotalPrice)
		if err != nil {
			return nil
		}

		customers = append(customers, &rpc.TopCustomer{
			CustomerId:          customer.CustomerID,
			CustomerName:        customer.CustomerName,
			NumerOfTransactions: int32(customer.TotalQuantity),
			SalesAmount:         totalPrice,
		})
	}

	err := stream.SendMsg(&rpc.StreamTopCustomersResponse{
		TopCustomers: customers,
	})
	if err != nil {
		fmt.Println("something went wrong while sending the message, err: ", err)
		return err
	}

	return nil
}

func streamSalesByProduct(productId string, salesByProduct *db.GetTotalSalesByProductIdRow, stream rpc.AnalticsService_StreamSalesByProductServer) error {
	amount, err := utils.PgNumericToFloat32(salesByProduct.TotalPrice)
	if err != nil {
		return err
	}

	// stream only if the request product id matches the created transaction product id
	if salesByProduct.ProductID == productId {
		err = stream.SendMsg(&rpc.StreamSalesByProductResponse{
			ProductId:     salesByProduct.ProductID,
			ProductName:   salesByProduct.ProductName,
			SalesQuantity: int32(salesByProduct.TotalQuantity),
			SalesAmount:   float32(amount),
		})

		if err != nil {
			fmt.Println("something went wrong while sending the message, err: ", err)
			return err
		}
	}

	return err
}

func streamTotalSales(totalSales *db.GetTotalSalesRow, stream rpc.AnalticsService_StreamTotalSalesServer) error {
	amount, err := utils.PgNumericToFloat32(totalSales.TotalPrice)
	if err != nil {
		return err
	}

	err = stream.SendMsg(&rpc.StreamTotalSalesResponse{
		TotalQuantity:        int32(totalSales.TotalQuantity),
		TotalAmount:          amount,
		NumberOfTransactions: int32(totalSales.TotalTransactions),
	})

	if err != nil {
		fmt.Println("something went wrong while sending the message, err: ", err)
		return err
	}

	return nil
}
