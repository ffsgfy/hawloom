-- name: CreateVordZero :exec
INSERT INTO vord (doc, num, flags, start_at, finish_at)
VALUES ($1, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
