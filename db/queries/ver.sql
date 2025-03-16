-- name: CreateVer :one
INSERT INTO ver (id, doc, vord_num, created_by, summary, content)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FindVer :one
SELECT sqlc.embed(ver), account.name AS author
FROM ver
    JOIN account ON account.id = ver.created_by
WHERE ver.id = $1;

-- name: FindVerWithVote :one
SELECT sqlc.embed(ver), account.name AS ver_author, doc.flags AS doc_flags,
    CAST(ver_vote.account IS NOT NULL AS BOOLEAN) AS ver_vote_exists,
    CAST(doc_vote.account IS NOT NULL AS BOOLEAN) AS doc_vote_exists
FROM ver
    JOIN account ON account.id = ver.created_by
    JOIN doc ON doc.id = ver.doc
    JOIN vord ON vord.doc = ver.doc AND vord.num = ver.vord_num
    LEFT JOIN vote AS ver_vote
        ON ver_vote.ver = $1 AND ver_vote.account = $2
    LEFT JOIN vote AS doc_vote
        ON doc_vote.doc = ver.doc AND doc_vote.vord_num = ver.vord_num AND doc_vote.account = $2
WHERE ver.id = $1
LIMIT 1;

-- name: FindVerForVote :one
SELECT ver.vord_num, ver.doc AS doc_id, doc.flags AS doc_flags,
    CAST(ver_vote.account IS NOT NULL AS BOOLEAN) AS ver_vote_exists,
    CAST(doc_vote.account IS NOT NULL AS BOOLEAN) AS doc_vote_exists
FROM ver
    JOIN doc ON doc.id = ver.doc
    JOIN vord ON vord.doc = ver.doc AND vord.num = ver.vord_num
    LEFT JOIN vote AS ver_vote
        ON ver_vote.ver = $1 AND ver_vote.account = $2
    LEFT JOIN vote AS doc_vote
        ON doc_vote.doc = ver.doc AND doc_vote.vord_num = ver.vord_num AND doc_vote.account = $2
WHERE ver.id = $1
LIMIT 1
FOR SHARE OF vord;

-- name: FindVerForDelete :one
SELECT ver.vord_num, ver.created_by, ver.doc AS doc_id
FROM ver
    JOIN vord ON vord.doc = ver.doc AND vord.num = ver.vord_num
WHERE ver.id = $1
FOR SHARE OF vord;

-- name: FindVersForCommit :many
-- Assumes vord is locked
SELECT id, votes FROM ver
WHERE doc = $1 AND vord_num = -1
ORDER BY votes DESC
LIMIT 2;

-- name: FindCurrentVer :one
SELECT sqlc.embed(ver), ver_acc.name AS ver_author,
    sqlc.embed(doc), doc_acc.name AS doc_author,
    sqlc.embed(vord)
FROM ver
    JOIN account AS ver_acc ON ver_acc.id = ver.created_by
    JOIN doc ON doc.id = ver.doc
    JOIN account AS doc_acc ON doc_acc.id = doc.created_by
    JOIN vord ON vord.doc = ver.doc AND vord.num = -1
WHERE ver.doc = $1
ORDER BY ver.vord_num DESC, ver.votes DESC
LIMIT 1;

-- name: FindWinningVer :one
SELECT sqlc.embed(ver), ver_acc.name AS ver_author,
    sqlc.embed(doc), doc_acc.name AS doc_author,
    sqlc.embed(vord)
FROM ver
    JOIN account AS ver_acc ON ver_acc.id = ver.created_by
    JOIN doc ON doc.id = ver.doc
    JOIN account AS doc_acc ON doc_acc.id = doc.created_by
    JOIN vord ON vord.doc = ver.doc AND vord.num = sqlc.arg(vord_num_join)
WHERE ver.doc = $1 AND ver.vord_num = $2
ORDER BY ver.votes DESC
LIMIT 1;

-- name: FindVerList :many
SELECT ver.id, ver.votes, account.name AS author, ver.summary
FROM ver
    JOIN account ON account.id = ver.created_by
WHERE ver.doc = $1 AND ver.vord_num = $2
ORDER BY ver.votes DESC;

-- name: UpdateVerVotes :exec
-- Assumes vord is locked
UPDATE ver
SET votes = votes + sqlc.arg(delta)
WHERE id = $1;

-- name: DeleteVer :exec
DELETE FROM ver WHERE id = $1;
