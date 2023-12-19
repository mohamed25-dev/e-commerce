-- name: GetTransactionById :one
SELECT * FROM transactions WHERE id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    customer_id,
    customer_name,
    product_id,
    product_name,
    quantity,
    total_price
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetTotalSales :one
SELECT 
	COUNT(id) as total_transactions,
	SUM(total_price) as total_price,
	SUM(quantity) as total_quantity
FROM transactions;

-- name: GetTotalSalesByProductId :one
SELECT 
  product_id,
  product_name,
	COUNT(id) as total_transactions,
	SUM(total_price) as total_price,
	SUM(quantity) as total_quantity
FROM transactions
WHERE product_id = $1
GROUP BY product_id, product_name;
