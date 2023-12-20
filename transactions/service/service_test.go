package service

import (
	"context"
	mockdb "ecommerce/transactions/db/mock"
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	query := mockdb.NewMockQuerier(ctrl)

	transactionRequest := models.CreateTransactionRequestModel{
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   4,
	}

	fetchedCustomer := db.Customer{
		ID: "123",
	}

	query.EXPECT().
		GetCustomerById(gomock.Any(), gomock.Eq(transactionRequest.CustomerID)).
		Times(1).
		Return(fetchedCustomer, nil)

	fetchedProduct := db.Product{
		ID: "123",
	}
	query.EXPECT().
		GetProductById(gomock.Any(), gomock.Eq(transactionRequest.ProductID)).
		Times(1).
		Return(fetchedProduct, nil)

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
