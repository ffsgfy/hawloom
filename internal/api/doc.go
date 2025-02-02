package api

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

type DocFlags int32

const (
	DocFlagPublic   DocFlags = 1 << 0
	DocFlagApproval DocFlags = 1 << 1
	DocFlagMajority DocFlags = 1 << 2
)

const (
	DocTitleMaxLength = 256
	VordMinDuration   = 10
)

type CreateDocParams struct {
	Title        string
	Summary      string
	Content      string
	Flags        DocFlags
	VordDuration int32
}

func (sc *StateCtx) CreateDoc(params *CreateDocParams) (*db.Doc, *db.Ver, error) {
	if len(params.Title) > DocTitleMaxLength {
		return nil, nil, ErrDocTitleTooLong
	}
	if params.VordDuration < VordMinDuration {
		return nil, nil, ErrRoundDurationTooSmall
	}

	authState, err := GetValidAuthState(sc.Ctx)
	if err != nil {
		return nil, nil, err
	}

	doc_id := uuid.New()
	ver_id := uuid.New()
	var doc *db.Doc
	var ver *db.Ver

	if err = sc.Tx(func(sc *StateCtx) error {
		if doc, err = sc.Queries.CreateDoc(sc.Ctx, &db.CreateDocParams{
			ID:           doc_id,
			Title:        params.Title,
			Flags:        int32(params.Flags),
			CreatedBy:    authState.Account.ID,
			VordDuration: params.VordDuration,
		}); err != nil {
			return err
		}

		if err = sc.Queries.CreateVordZero(sc.Ctx, doc_id); err != nil {
			return err
		}

		ver, err = sc.Queries.CreateVer(sc.Ctx, &db.CreateVerParams{
			ID:        ver_id,
			Doc:       doc_id,
			VordNum:   0,
			CreatedBy: authState.Account.ID,
			Summary:   params.Summary,
			Content:   params.Content,
		})

		return err
	}); err != nil {
		return nil, nil, err
	}

	// TODO: log flags
	ctxlog.Info(
		sc.Ctx, "doc created",
		"account_id", authState.Account.ID,
		"doc_id", doc_id,
		"ver_id", ver_id,
	)

	return doc, ver, nil
}

func (sc *StateCtx) DeleteDoc(id uuid.UUID) error {
	authState, err := GetValidAuthState(sc.Ctx)
	if err != nil {
		return err
	}

	doc, err := sc.Queries.FindDoc(sc.Ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDocNotFound
		}
		return err
	}

	if doc.CreatedBy != authState.Account.ID {
		return ErrForbidden
	}

	err = sc.Queries.DeleteDoc(sc.Ctx, id)
	if err != nil {
		return err
	}

	ctxlog.Info(
		sc.Ctx, "doc deleted",
		"account_id", authState.Account.ID,
		"doc_id", id,
	)

	return nil
}
