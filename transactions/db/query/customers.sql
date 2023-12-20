-- name: GetCustomerById :one
SELECT * FROM customers WHERE id = $1;

-- name: CreateCustomer :one
INSERT INTO customers (
    id,
    customer_name
) VALUES (
  $1, $2
) RETURNING *;