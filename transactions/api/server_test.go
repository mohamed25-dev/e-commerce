package api

import (
	"context"
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"ecommerce/transactions/proto"
	"ecommerce/transactions/utils"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionsService struct {
	mock.Mock
}

func (m *MockTransactionsService) GetTransactionById(ctx context.Context, id string) (db.Transaction, error) {

	val, err := utils.Float32ToPgNumeric(21)
	if err != nil {
		return db.Transaction{}, err
	}

	var expectedTransaction = db.Transaction{
		ID:         "123",
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   2,
		TotalPrice: val,
	}

	return expectedTransaction, nil
}

func (service *MockTransactionsService) CreateTransaction(ctx context.Context, transactionData models.CreateTransactionRequestModel) (db.Transaction, error) {
	return db.Transaction{}, nil
}

func TestGetTransactionById(t *testing.T) {

	request := &proto.GetTransactionByIdRequest{
		Id: "123",
	}

	mockService := new(MockTransactionsService)

	// ctrl := gomock.NewController(t)
	// query := mockdb.NewMockQuerier(ctrl)

	// query.EXPECT().
	// 	GetTransactionById(gomock.Any(), expectedTransaction.ID).
	// 	Times(1).
	// 	Return(expectedTransaction, nil)
	mockService.On("GetTransactionById", context.Background(), request.Id).Return("123", nil)

	server := &TransactionsServer{
		service: mockService,
		mu:      sync.Mutex{},
		streams: make(map[chan *proto.CreateTransactionResponse]struct{}),
	}

	// server.service = mockService

	response, err := server.GetTransactionById(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	mockService.AssertCalled(t, "GetTransactionById", mock.Anything, "your_transaction_id")

	mockService.AssertExpectations(t)
}
