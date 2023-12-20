package models

import (
	db "ecommerce/transactions/db/sqlc"
	"ecommerce/transactions/proto"
	"ecommerce/transactions/utils"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapRpcRequestToCreateTransactionModel(req *proto.CreateTransactionRequest) (CreateTransactionRequestModel, error) {
	data := CreateTransactionRequestModel{
		CustomerID: req.CustomerId,
		ProductID:  req.ProductId,
		Quantity:   req.Quantity,
	}

	return data, nil
}

func MapTransactionDataToAnalyticsTransaction(transaction CreateTransactionRequestModel, product db.Product, customer db.Customer) (CreateAnalyticsTransactionRequestModel, error) {
	data := CreateAnalyticsTransactionRequestModel{
		CustomerID:   transaction.CustomerID,
		CustomerName: customer.CustomerName,
		ProductID:    transaction.ProductID,
		ProductName:  product.ProductName,
		Quantity:     transaction.Quantity,
		TotalPrice:   transaction.TotalPrice,
	}

	fmt.Println(data)
	return data, nil
}

func MapTransactionDataToDbParams(data CreateTransactionRequestModel) (db.CreateTransactionParams, error) {
	totalPrice, err := utils.Float32ToPgNumeric(data.TotalPrice)
	if err != nil {
		return db.CreateTransactionParams{}, err
	}

	transactionData := db.CreateTransactionParams{
		CustomerID: data.CustomerID,
		ProductID:  data.ProductID,
		Quantity:   data.Quantity,
		TotalPrice: totalPrice,
	}

	return transactionData, nil
}

func MapDbTransactionToRpcResponse(trx db.Transaction) (*proto.GetTransactionByIdResponse, error) {
	amountString := trx.TotalPrice.Int.String()
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		fmt.Println("Something went wrong while converting string to number, err: ", err)
		return nil, err
	}

	var pgTimestamp pgtype.Timestamptz
	if err := trx.CreatedAt.ScanTimestamptz(pgTimestamp); err != nil {
		fmt.Println("Error scanning created_at, err: ", err)
		return nil, err
	}

	response := &proto.GetTransactionByIdResponse{
		Transaction: &proto.Transaction{
			Id:          trx.ID,
			CustomerId:  trx.CustomerID,
			ProductId:   trx.ProductID,
			Quantity:    trx.Quantity,
			TotalAmount: float32(amount),
			CreatedAt:   timestamppb.New(pgTimestamp.Time),
		},
	}

	return response, nil
}
