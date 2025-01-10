-- name: CreateVer :one
INSERT INTO ver (id, doc, vord_num, created_by, summary, content)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: FindVerForDelete :one
SELECT ver.vord_num, ver.created_by, ver.doc AS doc_id FROM ver
WHERE ver.id = $1
FOR UPDATE;

-- name: FindVerForVote :one
SELECT ver.vord_num, ver.created_by, ver.doc AS doc_id, doc.flags AS doc_flags,
    CAST(ver_vote.account IS NOT NULL AS BOOLEAN) AS ver_vote_exists,
    CAST(doc_vote.account IS NOT NULL AS BOOLEAN) AS doc_vote_exists
FROM ver
    JOIN doc ON doc.id = ver.doc
    LEFT JOIN vote AS ver_vote
        ON ver_vote.ver = $1 AND ver_vote.account = $2
    LEFT JOIN vote AS doc_vote
        ON doc_vote.doc = ver.doc AND doc_vote.vord_num = ver.vord_num AND doc_vote.account = $2
WHERE ver.id = $1
LIMIT 1
FOR SHARE OF ver;

-- name: DeleteVer :exec
DELETE FROM ver WHERE id = $1;
