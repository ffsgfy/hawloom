-- name: CreateVote :exec
INSERT INTO vote (ver, doc, vord_num, account)
VALUES ($1, $2, $3, $4);

-- name: FindVoteForDelete :one
SELECT vote.* FROM vote
    JOIN vord ON vord.doc = vote.doc AND vord.num = vote.vord_num
WHERE vote.ver = $1 AND vote.account = $2
FOR UPDATE OF vote
FOR SHARE OF vord;

-- name: CountVotes :many
SELECT ver, COUNT(*) AS votes FROM vote
WHERE doc = $1 AND vote.vord_num = $2
GROUP BY ver;

-- name: CountVoters :one
-- Assumes vord is locked
SELECT COUNT(DISTINCT account) AS voters FROM vote
WHERE doc = $1 AND vord_num = $2;

-- name: DeleteVote :exec
DELETE FROM vote WHERE ver = $1 AND account = $2;
