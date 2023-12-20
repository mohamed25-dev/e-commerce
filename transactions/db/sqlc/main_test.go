package db

import (
	"context"
	"ecommerce/transactions/utils"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/client"
)

type MockTemporalClient struct {
	mock.Mock
}

func (m *MockTemporalClient) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (*client.WorkflowRun, error) {
	argsMock := m.Called(ctx, options, workflow, args)
	return argsMock.Get(0).(*client.WorkflowRun), argsMock.Error(1)
}

func (m *MockTemporalClient) CancelWorkflow(ctx context.Context, s string, b string) error {
	return nil
}

var testQueries *Queries
var mockTemporalClient *MockTemporalClient

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../transactions/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr, err := utils.GetDbConnectionString()
	fmt.Println(connStr)
	if err != nil {
		log.Fatal("could not get DB connection string, err: ", err)
	}

	// Create a connection pool
	config, err := pgx.ParseConfig(connStr)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close(context.Background())

	testQueries = New(pool)
	mockTemporalClient = &MockTemporalClient{}

	os.Exit(m.Run())
}
