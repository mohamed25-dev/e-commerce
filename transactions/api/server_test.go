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
	args := m.Called(ctx, id)
	return args.Get(0).(db.Transaction), args.Error(1)
}

func (m *MockTransactionsService) CreateTransaction(ctx context.Context, transactionData models.CreateTransactionRequestModel) (db.Transaction, error) {
	args := m.Called(ctx, transactionData)
	return args.Get(0).(db.Transaction), args.Error(1)
}

func TestGetTransactionById(t *testing.T) {
	request := &proto.GetTransactionByIdRequest{
		Id: "123",
	}

	val, err := utils.Float32ToPgNumeric(21)
	if err != nil {
		t.Fatal(err)
	}

	var expectedTransaction = db.Transaction{
		ID:         "123",
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   2,
		TotalPrice: val,
	}
	mockService := new(MockTransactionsService)

	mockService.On("GetTransactionById", mock.Anything, request.Id).Return(expectedTransaction, nil)

	server := &TransactionsServer{
		service: mockService,
		mu:      sync.Mutex{},
		streams: make(map[chan *proto.CreateTransactionResponse]struct{}),
	}

	response, err := server.GetTransactionById(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)

	mockService.AssertCalled(t, "GetTransactionById", mock.Anything, request.Id)

	mockService.AssertExpectations(t)
}
