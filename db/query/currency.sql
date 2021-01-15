-- name: CreateCurrency :one
INSERT INTO currencies (
    code
) VALUES (
    $1
) RETURNING *;

-- name: GetCurrencyByID :one
SELECT * FROM currencies
WHERE id = $1 LIMIT 1;

-- name: GetCurrencyByCode :one
SELECT * FROM currencies
WHERE code = $1 LIMIT 1;