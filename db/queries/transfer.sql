-- name: CreateTransfer :one
INSERT INTO transfers (
  "from",
  "to",
  amount,
  created_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetTransfer :many
SELECT * FROM transfers
WHERE "from" = $1 OR "to" = $1;