// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: vord.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const commitVord = `-- name: CommitVord :exec
UPDATE vord AS v1
SET flags = $2,
    num = (
        SELECT MAX(num) + 1 FROM vord AS v2
        WHERE v2.doc = v1.doc
    )
WHERE v1.doc = $1 AND v1.num = -1
`

func (q *Queries) CommitVord(ctx context.Context, doc uuid.UUID, flags int32) error {
	_, err := q.db.Exec(ctx, commitVord, doc, flags)
	return err
}

const createVord = `-- name: CreateVord :execrows
INSERT INTO vord (doc, num, flags, start_at, finish_at)
VALUES ($1, -1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + CAST($2 AS INTERVAL))
ON CONFLICT DO NOTHING
`

func (q *Queries) CreateVord(ctx context.Context, doc uuid.UUID, duration time.Duration) (int64, error) {
	result, err := q.db.Exec(ctx, createVord, doc, duration)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const createVordZero = `-- name: CreateVordZero :exec
INSERT INTO vord (doc, num, flags, start_at, finish_at)
VALUES ($1, 0, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
`

func (q *Queries) CreateVordZero(ctx context.Context, doc uuid.UUID) error {
	_, err := q.db.Exec(ctx, createVordZero, doc)
	return err
}

const findVordForCommit = `-- name: FindVordForCommit :one
SELECT vord.doc, vord.num, vord.flags, vord.start_at, vord.finish_at, doc.id, doc.title, doc.description, doc.flags, doc.created_by, doc.created_at, doc.vord_duration FROM vord
    JOIN doc ON doc.id = vord.doc
WHERE num = -1 AND finish_at <= CURRENT_TIMESTAMP
ORDER BY finish_at
LIMIT 1
FOR UPDATE SKIP LOCKED
`

type FindVordForCommitRow struct {
	Vord Vord `db:"vord"`
	Doc  Doc  `db:"doc"`
}

func (q *Queries) FindVordForCommit(ctx context.Context) (*FindVordForCommitRow, error) {
	row := q.db.QueryRow(ctx, findVordForCommit)
	var i FindVordForCommitRow
	err := row.Scan(
		&i.Vord.Doc,
		&i.Vord.Num,
		&i.Vord.Flags,
		&i.Vord.StartAt,
		&i.Vord.FinishAt,
		&i.Doc.ID,
		&i.Doc.Title,
		&i.Doc.Description,
		&i.Doc.Flags,
		&i.Doc.CreatedBy,
		&i.Doc.CreatedAt,
		&i.Doc.VordDuration,
	)
	return &i, err
}

const findVordForCommitByDocID = `-- name: FindVordForCommitByDocID :one
SELECT vord.doc, vord.num, vord.flags, vord.start_at, vord.finish_at, doc.id, doc.title, doc.description, doc.flags, doc.created_by, doc.created_at, doc.vord_duration FROM vord
    JOIN doc ON doc.id = vord.doc
WHERE doc = $1 AND num = -1
FOR UPDATE NOWAIT
`

type FindVordForCommitByDocIDRow struct {
	Vord Vord `db:"vord"`
	Doc  Doc  `db:"doc"`
}

func (q *Queries) FindVordForCommitByDocID(ctx context.Context, doc uuid.UUID) (*FindVordForCommitByDocIDRow, error) {
	row := q.db.QueryRow(ctx, findVordForCommitByDocID, doc)
	var i FindVordForCommitByDocIDRow
	err := row.Scan(
		&i.Vord.Doc,
		&i.Vord.Num,
		&i.Vord.Flags,
		&i.Vord.StartAt,
		&i.Vord.FinishAt,
		&i.Doc.ID,
		&i.Doc.Title,
		&i.Doc.Description,
		&i.Doc.Flags,
		&i.Doc.CreatedBy,
		&i.Doc.CreatedAt,
		&i.Doc.VordDuration,
	)
	return &i, err
}

const lockVord = `-- name: LockVord :execrows
SELECT 1 FROM vord
WHERE doc = $1 AND num = -1
FOR SHARE
`

func (q *Queries) LockVord(ctx context.Context, doc uuid.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, lockVord, doc)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const updateVord = `-- name: UpdateVord :exec
UPDATE vord
SET flags = $2, finish_at = $3
WHERE doc = $1 AND num = -1
`

type UpdateVordParams struct {
	Doc      uuid.UUID `db:"doc"`
	Flags    int32     `db:"flags"`
	FinishAt time.Time `db:"finish_at"`
}

func (q *Queries) UpdateVord(ctx context.Context, arg *UpdateVordParams) error {
	_, err := q.db.Exec(ctx, updateVord, arg.Doc, arg.Flags, arg.FinishAt)
	return err
}
