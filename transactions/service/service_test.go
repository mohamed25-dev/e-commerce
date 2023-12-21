package service

import (
	"context"
	mockdb "ecommerce/transactions/db/mock"
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"ecommerce/transactions/utils"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/mocks"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	query := mockdb.NewMockQuerier(ctrl)

	tCilent := &mocks.Client{}
	transactionsService := TransactionsService{Queries: query, TemporalClient: tCilent}

	transactionRequest := models.CreateTransactionRequestModel{
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   4,
		TotalPrice: 60,
	}

	fetchedCustomer := db.Customer{
		ID: "123",
	}

	price, err := utils.Float32ToPgNumeric(15)
	if err != nil {
		return
	}
	fetchedProduct := db.Product{
		ID:    "123",
		Price: price,
	}

	// mocking queries
	query.EXPECT().
		GetCustomerById(gomock.Any(), "123").
		Times(1).
		Return(fetchedCustomer, nil)

	query.EXPECT().
		GetProductById(gomock.Any(), gomock.Eq(transactionRequest.ProductID)).
		Times(1).
		Return(fetchedProduct, nil)

	// mocking workflow
	wfRun := &mocks.WorkflowRun{}
	options := client.StartWorkflowOptions{
		ID:        "transactions-queue",
		TaskQueue: "transactions-queue",
	}

	tCilent.On("ExecuteWorkflow", context.Background(), options, mock.AnythingOfType("func(internal.Context, models.CreateTransactionRequestModel, db.Product, db.Customer) (db.Transaction, error)"),
		transactionRequest, fetchedProduct, fetchedCustomer).Return(wfRun, nil)

	createdTransaction := db.Transaction{}
	wfRun.On("Get", mock.Anything, &createdTransaction).Return(nil)

	// calling the function
	res, err := transactionsService.CreateTransaction(context.Background(), transactionRequest)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Empty(t, res)
}

func TestGetTransactionById(t *testing.T) {
	transaction := db.Transaction{
		ID: "123",
	}

	testCases := []struct {
		name          string
		transactionId string
		buildStubs    func(queries *mockdb.MockQuerier)
		checkResponse func(t *testing.T, transaction db.Transaction, err error)
	}{
		{
			name:          "success",
			transactionId: "123",
			buildStubs: func(query *mockdb.MockQuerier) {
				query.EXPECT().
					GetTransactionById(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, nil)
			},
			checkResponse: func(t *testing.T, receivedTransaction db.Transaction, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, receivedTransaction)
				require.Equal(t, transaction.ID, receivedTransaction.ID)
			},
		},
		{
			name:          "failed",
			transactionId: "123",
			buildStubs: func(query *mockdb.MockQuerier) {
				query.EXPECT().
					GetTransactionById(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(db.Transaction{}, errors.New("something went wrong"))
			},
			checkResponse: func(t *testing.T, receivedTransaction db.Transaction, err error) {
				require.Error(t, err)
				require.Empty(t, receivedTransaction.ID)
			},
		},
	}

	for _, testCase := range testCases {
		ctrl := gomock.NewController(t)
		query := mockdb.NewMockQuerier(ctrl)

		t.Run(testCase.name, func(t *testing.T) {
			testCase.buildStubs(query)
			service := TransactionsService{Queries: query}
			returnedTransaction, err := service.GetTransactionById(context.Background(), "123")

			testCase.checkResponse(t, returnedTransaction, err)
		})
	}
}
