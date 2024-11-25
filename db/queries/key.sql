-- name: FindKeys :many
SELECT * FROM key;

-- name: CreateKey :one
INSERT INTO key (data) VALUES ($1) RETURNING *;
