-- name: GetCustomerById :one
SELECT * FROM customers WHERE id = $1;