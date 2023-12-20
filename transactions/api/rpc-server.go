package server

import (
	"context"
	"ecommerce/transactions/models"
	"ecommerce/transactions/proto"
	"ecommerce/transactions/service"
	"ecommerce/transactions/utils"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func InitTransactionsRpcServer(service *service.TransactionsService, streams *map[chan *proto.CreateTransactionResponse]struct{}) *grpc.Server {
	grpcServer := grpc.NewServer()
	proto.RegisterTransactionServiceServer(grpcServer, &TransactionsServer{service: service, streams: *streams})

	return grpcServer
}

type TransactionsServer struct {
	proto.UnimplementedTransactionServiceServer
	service *service.TransactionsService
	mu      sync.Mutex
	streams map[chan *proto.CreateTransactionResponse]struct{}
}

func (s *TransactionsServer) GetTransactionById(ctx context.Context, req *proto.GetTransactionByIdRequest) (*proto.GetTransactionByIdResponse, error) {
	trx, err := s.service.GetTransactionById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	response, err := models.MapDbTransactionToRpcResponse(trx)
	if err != nil {
		return nil, nil
	}

	return response, nil
}

func (s *TransactionsServer) CreateTransaction(ctx context.Context, req *proto.CreateTransactionRequest) (*proto.CreateTransactionResponse, error) {
	transactionData := models.CreateTransactionData{
		CustomerID: req.CustomerId,
		ProductID:  req.ProductId,
		Quantity:   req.Quantity,
	}

	createdTransaction, err := s.service.CreateTransaction(ctx, transactionData)
	if err != nil {
		fmt.Println("something went wrong while creating the transaction, err: ", err)
		return nil, err
	}

	amount, err := utils.PgNumericToFloat32(createdTransaction.TotalPrice)
	if err != nil {
		fmt.Println("Something went wrong while converting string to number, err: ", err)
		return nil, err
	}

	var pgTimestamp pgtype.Timestamptz
	if err := createdTransaction.CreatedAt.ScanTimestamptz(pgTimestamp); err != nil {
		fmt.Println("Something went wrong while converting pgtype.Timestamps to time.Timestamps, err: ", err)
		return nil, err
	}

	response := &proto.CreateTransactionResponse{
		Transaction: &proto.Transaction{
			Id:          createdTransaction.ID,
			CustomerId:  createdTransaction.CustomerID,
			ProductId:   createdTransaction.ProductID,
			Quantity:    createdTransaction.Quantity,
			TotalAmount: amount,
			CreatedAt:   timestamppb.New(pgTimestamp.Time),
		},
	}

	s.mu.Lock()
	for ch := range s.streams {
		ch <- response
	}
	s.mu.Unlock()

	return response, nil
}

func (s *TransactionsServer) StreamTransactions(empty *proto.Empty, stream proto.TransactionService_StreamTransactionsServer) error {

	s.mu.Lock()
	ch := make(chan *proto.CreateTransactionResponse)
	s.streams[ch] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.streams, ch)
		close(ch)
		s.mu.Unlock()
	}()

	for {
		transaction, ok := <-ch
		if !ok {
			fmt.Println("channel is no longer available !!")
			break
		}

		err := stream.SendMsg(&proto.StreamTransactionResponse{Transaction: &proto.Transaction{Id: transaction.Transaction.Id, TotalAmount: transaction.Transaction.TotalAmount}})
		if err != nil {
			fmt.Println("something went wrong while sending the message, err: ", err)
			return err
		}

	}

	return nil
}
