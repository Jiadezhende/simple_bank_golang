-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  transfer_id,
  amount,
  entry_type,
  currency,
  created_at
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;