-- name: CreateTransfer :one
INSERT INTO transfers (
    user_id,
    currency_id,
    amount,
    time_placed
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE
    user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;