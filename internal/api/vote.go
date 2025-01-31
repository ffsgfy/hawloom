package api

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func (sc *StateCtx) CreateVote(verID uuid.UUID) error {
	authState, err := GetValidAuthState(sc.Ctx)
	if err != nil {
		return err
	}

	var docID uuid.UUID
	var approval bool

	if err = sc.Tx(func(sc *StateCtx) error {
		row, err := sc.Queries.FindVerForVote(sc.Ctx, verID, authState.Account.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrVerNotFound
			}
			return err
		}
		docID = row.DocID
		approval = utils.TestFlags(row.DocFlags, DocFlagApproval)

		if row.VordNum != -1 {
			return ErrVotePastVer
		}
		if row.CreatedBy != nil && *row.CreatedBy == authState.Account.ID {
			return ErrVoteOwnVer
		}
		if row.VerVoteExists {
			return ErrVerVoteExists
		}

		if approval {
			if row.DocVoteExists {
				return ErrDocVoteExists
			}
		}

		if err = sc.Queries.CreateVote(sc.Ctx, &db.CreateVoteParams{
			Ver:     verID,
			Doc:     row.DocID,
			VordNum: row.VordNum,
			Account: authState.Account.ID,
		}); err != nil {
			return err
		}

		return sc.Queries.UpdateVerVotes(sc.Ctx, verID, 1)
	}); err != nil {
		return err
	}

	ctxlog.Info(
		sc.Ctx, "vote created",
		"account_id", authState.Account.ID,
		"doc_id", docID,
		"ver_id", verID,
		"approval", approval,
	)

	return nil
}

func (sc *StateCtx) DeleteVote(verID uuid.UUID) error {
	authState, err := GetValidAuthState(sc.Ctx)
	if err != nil {
		return err
	}

	var vote *db.Vote

	if err = sc.Tx(func(sc *StateCtx) error {
		if vote, err = sc.Queries.FindVoteForDelete(
			sc.Ctx, verID, authState.Account.ID,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrVoteNotFound
			}
			return err
		}

		if vote.VordNum != -1 {
			return ErrDeletePastVote
		}

		if err = sc.Queries.DeleteVote(
			sc.Ctx, verID, authState.Account.ID,
		); err != nil {
			return err
		}

		return sc.Queries.UpdateVerVotes(sc.Ctx, verID, -1)
	}); err != nil {
		return err
	}

	ctxlog.Info(
		sc.Ctx, "vote deleted",
		"account_id", authState.Account.ID,
		"doc_id", vote.Doc,
		"ver_id", vote.Ver,
	)

	return nil
}
