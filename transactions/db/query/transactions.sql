-- name: GetTransactionById :one
SELECT * FROM transactions WHERE id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (
    id,
    customer_id,
    product_id,
    quantity,
    total_price
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;