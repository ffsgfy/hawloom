// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: vote.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const countVoters = `-- name: CountVoters :one
SELECT COUNT(DISTINCT account) AS voters FROM vote
WHERE doc = $1 AND vord_num = $2
`

type CountVotersParams struct {
	Doc     uuid.UUID `db:"doc"`
	VordNum int32     `db:"vord_num"`
}

// Assumes vord is locked
func (q *Queries) CountVoters(ctx context.Context, arg *CountVotersParams) (int64, error) {
	row := q.db.QueryRow(ctx, countVoters, arg.Doc, arg.VordNum)
	var voters int64
	err := row.Scan(&voters)
	return voters, err
}

const createVote = `-- name: CreateVote :exec
INSERT INTO vote (ver, doc, vord_num, account)
VALUES ($1, $2, $3, $4)
`

type CreateVoteParams struct {
	Ver     uuid.UUID `db:"ver"`
	Doc     uuid.UUID `db:"doc"`
	VordNum int32     `db:"vord_num"`
	Account int32     `db:"account"`
}

func (q *Queries) CreateVote(ctx context.Context, arg *CreateVoteParams) error {
	_, err := q.db.Exec(ctx, createVote,
		arg.Ver,
		arg.Doc,
		arg.VordNum,
		arg.Account,
	)
	return err
}

const deleteVote = `-- name: DeleteVote :exec
DELETE FROM vote WHERE ver = $1 AND account = $2
`

type DeleteVoteParams struct {
	Ver     uuid.UUID `db:"ver"`
	Account int32     `db:"account"`
}

func (q *Queries) DeleteVote(ctx context.Context, arg *DeleteVoteParams) error {
	_, err := q.db.Exec(ctx, deleteVote, arg.Ver, arg.Account)
	return err
}

const findVoteForDelete = `-- name: FindVoteForDelete :one
SELECT vote.ver, vote.doc, vote.vord_num, vote.account FROM vote
    JOIN vord ON vord.doc = vote.doc AND vord.num = vote.vord_num
WHERE vote.ver = $1 AND vote.account = $2
FOR UPDATE OF vote
FOR SHARE OF vord
`

type FindVoteForDeleteParams struct {
	Ver     uuid.UUID `db:"ver"`
	Account int32     `db:"account"`
}

func (q *Queries) FindVoteForDelete(ctx context.Context, arg *FindVoteForDeleteParams) (*Vote, error) {
	row := q.db.QueryRow(ctx, findVoteForDelete, arg.Ver, arg.Account)
	var i Vote
	err := row.Scan(
		&i.Ver,
		&i.Doc,
		&i.VordNum,
		&i.Account,
	)
	return &i, err
}
