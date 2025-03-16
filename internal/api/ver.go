package api

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

type CreateVerParams struct {
	DocID   uuid.UUID
	Summary string
	Content string
}

func (sc *StateCtx) CreateVer(params *CreateVerParams) (*db.Ver, error) {
	authToken, err := GetValidAuthToken(sc.Ctx)
	if err != nil {
		return nil, err
	}

	// TODO: limit max summary/content size

	var ver *db.Ver

	if err = sc.Tx(func(sc *StateCtx) error {
		if _, err := sc.Queries.LockVord(sc.Ctx, params.DocID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNoVordExists
			}
			return err
		}

		ver, err = sc.Queries.CreateVer(sc.Ctx, &db.CreateVerParams{
			ID:        uuid.New(),
			Doc:       params.DocID,
			VordNum:   -1,
			CreatedBy: authToken.AccountID,
			Summary:   params.Summary,
			Content:   params.Content,
		})

		return err
	}); err != nil {
		return nil, err
	}

	ctxlog.Info(
		sc.Ctx, "ver created",
		"account_id", authToken.AccountID,
		"doc_id", ver.Doc,
		"ver_id", ver.ID,
	)

	return ver, nil
}

func (sc *StateCtx) DeleteVer(id uuid.UUID) (uuid.UUID, error) {
	authToken, err := GetValidAuthToken(sc.Ctx)
	if err != nil {
		return uuid.UUID{}, err
	}

	var docID uuid.UUID

	if err = sc.Tx(func(sc *StateCtx) error {
		row, err := sc.Queries.FindVerForDelete(sc.Ctx, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrVerNotFound
			}
			return err
		}
		docID = row.DocID

		if row.CreatedBy != authToken.AccountID {
			return ErrForbidden
		}

		if row.VordNum != -1 {
			return ErrDeletePastVer
		}

		return sc.Queries.DeleteVer(sc.Ctx, id)
	}); err != nil {
		return uuid.UUID{}, err
	}

	ctxlog.Info(
		sc.Ctx, "ver deleted",
		"account_id", authToken.AccountID,
		"doc_id", docID,
		"ver_id", id,
	)

	return docID, nil
}
