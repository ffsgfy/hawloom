-- name: CreateDoc :one
INSERT INTO doc (id, title, flags, created_by, vord_duration)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: FindDoc :one
SELECT * FROM doc WHERE id = $1;

-- name: DeleteDoc :exec
DELETE FROM doc WHERE id = $1;
