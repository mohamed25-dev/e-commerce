package workflow

import (
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/models"
	"ecommerce/transactions/utils"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func TestCreateTransactionWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	activities := TransactionActivity{}
	wf := TransactionWorkflow{}

	env.RegisterActivity(activities.SendTransactionToAnalyticsActivity)
	env.RegisterActivity(activities.CreateTransactionActivity)

	val, err := utils.Float32ToPgNumeric(5)
	if err != nil {
		t.Fatal(err)
	}

	product := db.Product{
		ID:          "123",
		ProductName: utils.GenerateRandomString(8),
		Price:       val,
	}
	customer := db.Customer{
		ID:           "123",
		CustomerName: utils.GenerateRandomString(8),
	}

	transactionData := models.CreateTransactionRequestModel{
		ID:         "123",
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   3,
		TotalPrice: 15,
	}

	analyticsTransactionData := models.CreateAnalyticsTransactionRequestModel{
		ProductID:    product.ID,
		ProductName:  product.ProductName,
		CustomerID:   customer.ID,
		CustomerName: customer.CustomerName,
		TotalPrice:   transactionData.TotalPrice,
		Quantity:     transactionData.Quantity,
	}

	val, err = utils.Float32ToPgNumeric(transactionData.TotalPrice)
	if err != nil {
		t.Fatal(err)
	}

	createdTransaction := db.Transaction{
		ID:         transactionData.ID,
		ProductID:  transactionData.ProductID,
		CustomerID: transactionData.CustomerID,
		TotalPrice: val,
		Quantity:   transactionData.Quantity,
	}

	env.OnActivity(activities.CreateTransactionActivity, mock.Anything, transactionData).Return(createdTransaction, nil)
	env.OnActivity(activities.SendTransactionToAnalyticsActivity, mock.Anything, analyticsTransactionData).Return(nil)

	env.ExecuteWorkflow(wf.CreateTransactionWorkflow, transactionData, product, customer)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func TestCreateTransactionWorkflow_TransactionInsertionFailed(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	activities := TransactionActivity{}
	wf := TransactionWorkflow{}

	env.RegisterActivity(activities.SendTransactionToAnalyticsActivity)
	env.RegisterActivity(activities.CreateTransactionActivity)

	val, err := utils.Float32ToPgNumeric(5)
	if err != nil {
		t.Fatal(err)
	}

	product := db.Product{
		ID:          "123",
		ProductName: utils.GenerateRandomString(8),
		Price:       val,
	}
	customer := db.Customer{
		ID:           "123",
		CustomerName: utils.GenerateRandomString(8),
	}

	transactionData := models.CreateTransactionRequestModel{
		ID:         "123",
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   3,
		TotalPrice: 15,
	}

	analyticsTransactionData := models.CreateAnalyticsTransactionRequestModel{
		ProductID:    product.ID,
		ProductName:  product.ProductName,
		CustomerID:   customer.ID,
		CustomerName: customer.CustomerName,
		TotalPrice:   transactionData.TotalPrice,
		Quantity:     transactionData.Quantity,
	}

	env.OnActivity(activities.CreateTransactionActivity, mock.Anything, transactionData).Return(db.Transaction{}, errors.New("error"))
	env.OnActivity(activities.SendTransactionToAnalyticsActivity, mock.Anything, analyticsTransactionData).Return(nil)

	env.ExecuteWorkflow(wf.CreateTransactionWorkflow, transactionData, product, customer)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

func TestCreateTransactionWorkflow_CreatingAnalyticsTransactionFailed(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	activities := TransactionActivity{}
	wf := TransactionWorkflow{}

	env.RegisterActivity(activities.SendTransactionToAnalyticsActivity)
	env.RegisterActivity(activities.CreateTransactionActivity)

	val, err := utils.Float32ToPgNumeric(5)
	if err != nil {
		t.Fatal(err)
	}

	product := db.Product{
		ID:          "123",
		ProductName: utils.GenerateRandomString(8),
		Price:       val,
	}
	customer := db.Customer{
		ID:           "123",
		CustomerName: utils.GenerateRandomString(8),
	}

	transactionData := models.CreateTransactionRequestModel{
		ID:         "123",
		CustomerID: "123",
		ProductID:  "123",
		Quantity:   3,
		TotalPrice: 15,
	}

	val, err = utils.Float32ToPgNumeric(transactionData.TotalPrice)
	if err != nil {
		t.Fatal(err)
	}

	createdTransaction := db.Transaction{
		ID:         transactionData.ID,
		ProductID:  transactionData.ProductID,
		CustomerID: transactionData.CustomerID,
		TotalPrice: val,
		Quantity:   transactionData.Quantity,
	}

	analyticsTransactionData := models.CreateAnalyticsTransactionRequestModel{
		ProductID:    product.ID,
		ProductName:  product.ProductName,
		CustomerID:   customer.ID,
		CustomerName: customer.CustomerName,
		TotalPrice:   transactionData.TotalPrice,
		Quantity:     transactionData.Quantity,
	}

	env.OnActivity(activities.CreateTransactionActivity, mock.Anything, transactionData).Return(createdTransaction, nil)
	env.OnActivity(activities.SendTransactionToAnalyticsActivity, mock.Anything, analyticsTransactionData).Return(errors.New("error"))

	env.ExecuteWorkflow(wf.CreateTransactionWorkflow, transactionData, product, customer)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}
