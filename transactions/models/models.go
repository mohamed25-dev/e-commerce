package models

type CreateTransactionData struct {
	ID         string
	CustomerID string
	ProductID  string
	TotalPrice float32
	Quantity   int32
}

type CreateAnalyticsTransactionData struct {
	CustomerID   string
	CustomerName string
	ProductID    string
	ProductName  string
	TotalPrice   float32
	Quantity     int32
}
