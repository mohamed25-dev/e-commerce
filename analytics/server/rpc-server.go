package server

import (
	"context"
	db "ecommerce/analytics/db/sqlc"
	rpc "ecommerce/analytics/proto"
	"ecommerce/analytics/service"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc"
)

func InitAnalyticsRpcServer(service *service.AnalyticsService, totalSalesStreams map[chan *db.GetTotalSalesRow]struct{}, topCustomerStreams map[chan string]struct{}, salesByProductStreams map[chan *db.GetTotalSalesByProductIdRow]struct{}) *grpc.Server {
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

type AnalyticsServer struct {
	rpc.UnimplementedAnalticsServiceServer
	service *service.AnalyticsService
	// NOTE: creating two different mutexes might be a good idea,
	// as it will reduce the blocked code area for each mutex
	mu                    sync.Mutex
	topCustomerStreams    map[chan string]struct{}
	totalSalesStreams     map[chan *db.GetTotalSalesRow]struct{}
	salesByProductStreams map[chan *db.GetTotalSalesByProductIdRow]struct{}
}

func (s *AnalyticsServer) CreateAnalyticsTransaction(ctx context.Context, req *rpc.CreateAnalyticsTransactionRequest) (*rpc.CreateAnaltyicsTransactionResponse, error) {
	var totalPrice pgtype.Numeric
	//TODO: fix conversion issue which happens here.
	if err := totalPrice.Scan(fmt.Sprintf("%.2f", req.TotalAmount)); err != nil {
		fmt.Println("scanning error is: ", err)
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

		//TODO: enhance the mapping by creating a separate package for it
		amountString := totalSales.TotalPrice.Int.String()
		amount, err := strconv.Atoi(amountString)
		if err != nil {
			log.Println("Something went wrong while converting string to number, err: ", err)
			return err
		}

		err = stream.SendMsg(&rpc.StreamTotalSalesResponse{
			TotalQuantity:        int32(totalSales.TotalQuantity),
			TotalAmount:          float32(amount),
			NumberOfTransactions: float32(totalSales.TotalTransactions),
		})

		if err != nil {
			fmt.Println("something went wrong while sending the message, err: ", err)
			return err
		}

	}

	return nil
}

func (s *AnalyticsServer) StreamSalesByProduct(req *rpc.StreamSalesByProductRequest, stream rpc.AnalticsService_StreamSalesByProductServer) error {
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

		//TODO: enhance the mapping by creating a separate package for it
		amountString := salesByProduct.TotalPrice.Int.String()
		amount, err := strconv.Atoi(amountString)
		if err != nil {
			log.Println("Something went wrong while converting string to number, err: ", err)
			return err
		}

		// stream only if the request product id matches the created transaction product id
		if salesByProduct.ProductID == req.ProductId {
			err = stream.SendMsg(&rpc.StreamSalesByProductResponse{
				ProductId:     salesByProduct.ProductID,
				ProductName:   "Pixel 8", //TODO: remove hardcoded value after adding missing db columns
				SalesQuantity: int32(salesByProduct.TotalQuantity),
				SalesAmount:   float32(amount),
			})

			if err != nil {
				fmt.Println("something went wrong while sending the message, err: ", err)
				return err
			}
		}

		fmt.Println("+==============")
		fmt.Println(req.ProductId)
		fmt.Println(salesByProduct.ProductID)
		fmt.Println("+==============")

	}

	return nil
}
