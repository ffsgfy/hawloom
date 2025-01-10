-- name: CreateVote :exec
INSERT INTO vote (ver, doc, vord_num, account)
VALUES ($1, $2, $3, $4);

-- name: FindVoteForDelete :one
SELECT * FROM vote
WHERE ver = $1 AND account = $2
FOR UPDATE;

-- name: DeleteVote :exec
DELETE FROM vote WHERE ver = $1 AND account = $2;
