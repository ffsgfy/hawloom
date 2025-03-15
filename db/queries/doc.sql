-- name: CreateDoc :one
INSERT INTO doc (id, title, description, flags, created_by, vord_duration)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FindDoc :one
SELECT * FROM doc WHERE id = $1;

-- name: FindDocList :many
SELECT * FROM doc
WHERE created_by = $1
ORDER BY created_at DESC;

-- name: FindPublicDocList :many
SELECT * FROM doc
WHERE created_by = $1
    AND doc.flags & 1 = 1 -- DocFlagPublic
ORDER BY created_at DESC;

-- name: FindAllPublicDocList :many
SELECT doc.*, account.name AS author
FROM doc
    JOIN account ON account.id = doc.created_by
WHERE doc.flags & 1 = 1 -- DocFlagPublic
ORDER BY doc.created_at DESC;

-- name: DeleteDoc :exec
DELETE FROM doc WHERE id = $1;
