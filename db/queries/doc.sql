-- name: CreateDoc :one
INSERT INTO doc (id, title, flags, created_by, vord_duration, current_ver)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FindDoc :one
SELECT * FROM doc WHERE id = $1;

-- name: DeleteDoc :exec
DELETE FROM doc WHERE id = $1;
