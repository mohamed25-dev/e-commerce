-- name: GetRecentAnalytic :one
SELECT * FROM analytics ORDER BY created_at DESC LIMIT 1;

-- name: CreateAnalytic :one
INSERT INTO analytics (
    id,
    top_customers,
    total_sales
) VALUES (
  $1, $2, $3
) RETURNING *;