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

type CreateDocParams struct {
	Title        string
	Summary      string
	Content      string
	Flags        DocFlags
	VordDuration int32
}

func (sc *StateCtx) CreateDoc(params *CreateDocParams) (*db.Doc, *db.Ver, error) {
	if len(params.Title) > sc.Config.Doc.TitleMaxLength.V {
		return nil, nil, ErrDocTitleTooLong
	}
	if params.VordDuration < sc.Config.Vord.MinDuration.V {
		return nil, nil, ErrRoundDurationTooSmall
	}

	// TODO: limit max summary/content size

	authToken, err := GetValidAuthToken(sc.Ctx)
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
			CreatedBy:    authToken.AccountID,
			VordDuration: params.VordDuration,
		}); err != nil {
			return err
		}

		if err = sc.Queries.CreateVordZero(sc.Ctx, doc_id); err != nil {
			return err
		}

		if ver, err = sc.Queries.CreateVer(sc.Ctx, &db.CreateVerParams{
			ID:        ver_id,
			Doc:       doc_id,
			VordNum:   0,
			CreatedBy: authToken.AccountID,
			Summary:   params.Summary,
			Content:   params.Content,
		}); err != nil {
			return err
		}

		return sc.CreateVord(doc_id, params.VordDuration)
	}); err != nil {
		return nil, nil, err
	}

	// TODO: log flags
	ctxlog.Info(
		sc.Ctx, "doc created",
		"account_id", authToken.AccountID,
		"doc_id", doc_id,
		"ver_id", ver_id,
	)

	return doc, ver, nil
}

func (sc *StateCtx) DeleteDoc(id uuid.UUID) error {
	authToken, err := GetValidAuthToken(sc.Ctx)
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

	if doc.CreatedBy != authToken.AccountID {
		return ErrForbidden
	}

	err = sc.Queries.DeleteDoc(sc.Ctx, id)
	if err != nil {
		return err
	}

	ctxlog.Info(
		sc.Ctx, "doc deleted",
		"account_id", authToken.AccountID,
		"doc_id", id,
	)

	return nil
}
