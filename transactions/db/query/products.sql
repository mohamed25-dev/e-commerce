-- name: GetProductById :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (
    id,
    product_name,
    price
) VALUES (
  $1, $2, $3
) RETURNING *;