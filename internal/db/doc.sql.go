// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: doc.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createDoc = `-- name: CreateDoc :one
INSERT INTO doc (id, title, flags, created_by, vord_duration, current_ver)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, title, flags, created_by, created_at, vord_duration, current_ver
`

type CreateDocParams struct {
	ID           uuid.UUID `db:"id"`
	Title        string    `db:"title"`
	Flags        int32     `db:"flags"`
	CreatedBy    int32     `db:"created_by"`
	VordDuration int32     `db:"vord_duration"`
	CurrentVer   uuid.UUID `db:"current_ver"`
}

func (q *Queries) CreateDoc(ctx context.Context, arg *CreateDocParams) (*Doc, error) {
	row := q.db.QueryRow(ctx, createDoc,
		arg.ID,
		arg.Title,
		arg.Flags,
		arg.CreatedBy,
		arg.VordDuration,
		arg.CurrentVer,
	)
	var i Doc
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Flags,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.VordDuration,
		&i.CurrentVer,
	)
	return &i, err
}

const deleteDoc = `-- name: DeleteDoc :exec
DELETE FROM doc WHERE id = $1
`

func (q *Queries) DeleteDoc(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteDoc, id)
	return err
}

const findDoc = `-- name: FindDoc :one
SELECT id, title, flags, created_by, created_at, vord_duration, current_ver FROM doc WHERE id = $1
`

func (q *Queries) FindDoc(ctx context.Context, id uuid.UUID) (*Doc, error) {
	row := q.db.QueryRow(ctx, findDoc, id)
	var i Doc
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Flags,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.VordDuration,
		&i.CurrentVer,
	)
	return &i, err
}
