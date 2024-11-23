-- name: CreateAccount :one
INSERT INTO account (name, password_hash)
VALUES ($1, $2)
ON CONFLICT (name) DO NOTHING
RETURNING id;

-- name: FindAccountByID :one
SELECT * FROM account WHERE id = $1;

-- name: FindAccountByName :one
SELECT * FROM account WHERE name = $1;

-- name: CheckAccountName :one
SELECT COUNT(1) FROM account WHERE name = $1;
