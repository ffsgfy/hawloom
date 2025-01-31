-- name: CreateVordZero :exec
INSERT INTO vord (doc, num, flags, start_at, finish_at)
VALUES ($1, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: CreateVord :execrows
INSERT INTO vord (doc, num, flags, start_at, finish_at)
VALUES ($1, -1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + CAST(sqlc.arg(duration) AS INTERVAL))
ON CONFLICT DO NOTHING;

-- name: LockVord :execrows
SELECT 1 FROM vord
WHERE doc = $1 AND num = -1
FOR SHARE;

-- name: FindVordForCommitByDocID :one
SELECT sqlc.embed(vord), sqlc.embed(doc) FROM vord
    JOIN doc ON doc.id = vord.doc
WHERE doc = $1 AND num = -1
FOR UPDATE NOWAIT;

-- name: FindVordForCommit :one
SELECT sqlc.embed(vord), sqlc.embed(doc) FROM vord
    JOIN doc ON doc.id = vord.doc
WHERE num = -1 AND finish_at <= CURRENT_TIMESTAMP
ORDER BY finish_at
LIMIT 1
FOR UPDATE SKIP LOCKED;

-- name: UpdateVord :exec
UPDATE vord
SET flags = $2, finish_at = $3
WHERE doc = $1 AND num = -1;

-- name: CommitVord :exec
UPDATE vord AS v1
SET flags = $2,
    num = (
        SELECT MAX(num) + 1 FROM vord AS v2
        WHERE v2.doc = v1.doc
    )
WHERE v1.doc = $1 AND v1.num = -1;
