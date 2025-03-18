package api

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

type VordFlags int32

const (
	VordFlagError      VordFlags = 1 << 0 // unspecified internal error
	VordFlagTie        VordFlags = 1 << 1 // more than one ver with max votes
	VordFlagNoMajority VordFlags = 1 << 2 // majority rule not satisfied
	VordFlagNoVotes    VordFlags = 1 << 3 // no votes at all
)

func (sc *StateCtx) CreateVord(docID uuid.UUID, vordDuration int32) error {
	res, err := sc.Queries.CreateVord(sc.Ctx, docID, time.Second*time.Duration(vordDuration))
	if err != nil {
		return err
	}
	if res == 0 {
		return ErrVordExists
	}
	return nil
}

// Assumes vord is already locked
func (sc *StateCtx) commitVord(vord *db.Vord, majorityRule bool) (VordFlags, error) {
	if vord.Num != -1 {
		return VordFlagError, ErrCommitPastVord
	}

	verVotes, err := sc.Queries.FindVersForCommit(sc.Ctx, vord.Doc)
	if err != nil {
		return VordFlagError, err
	}
	if len(verVotes) == 0 || verVotes[0].Votes == 0 {
		return VordFlagNoVotes, nil
	}
	if len(verVotes) > 1 && verVotes[0].Votes == verVotes[1].Votes {
		return VordFlagTie, nil
	}

	if majorityRule {
		voters, err := sc.Queries.CountVoters(sc.Ctx, vord.Doc, vord.Num)
		if err != nil {
			return VordFlagError, err
		}
		if int64(verVotes[0].Votes)*2 <= voters {
			return VordFlagNoMajority, nil
		}
	}

	if err = sc.Queries.CommitVord(sc.Ctx, vord.Doc, 0); err != nil {
		return VordFlagError, err
	}

	return 0, nil
}

func (sc *StateCtx) findAndCommitVord() error {
	return sc.TxWith(pgx.TxOptions{IsoLevel: pgx.ReadCommitted}, func(sc *StateCtx) error {
		row, err := sc.Queries.FindVordForCommit(sc.Ctx)
		if err != nil {
			return err // may be sql.ErrNoRows
		}

		ctxlog.Info(sc.Ctx, "autocommit: committing vord", "doc_id", row.Doc.ID)
		flags, err := sc.commitVord(&row.Vord, utils.TestFlags(row.Doc.Flags, DocFlagMajority))

		var delay time.Duration
		if err != nil {
			delay = time.Second * time.Duration(sc.Config.Vord.MinDuration.V)
		} else if flags != 0 {
			delay = time.Second * time.Duration(max(
				float64(row.Doc.VordDuration)*sc.Config.Vord.DurationExtension.V,
				float64(sc.Config.Vord.MinDuration.V),
			))
		}

		if err != nil || flags != 0 {
			loglvl := ctxlog.INFO
			if err != nil {
				loglvl = ctxlog.ERROR
			}
			ctxlog.Log(
				sc.Ctx, loglvl, "autocommit: postponing vord commit",
				"doc_id", row.Doc.ID, "err", err, "flags", flags, "delay", delay,
			)

			if err = errors.Join(err, sc.Queries.UpdateVord(sc.Ctx, &db.UpdateVordParams{
				Doc:      row.Doc.ID,
				Flags:    int32(flags),
				FinishAt: time.Now().Add(delay),
			})); err != nil {
				return err
			}
			return nil
		}

		ctxlog.Info(sc.Ctx, "autocommit: creating new vord", "doc_id", row.Doc.ID)
		return sc.CreateVord(row.Doc.ID, row.Doc.VordDuration)
	})
}

func (sc *StateCtx) RunAutocommit() {
	sc.TasksWG.Add(1)
	defer sc.TasksWG.Done()

	for {
		if sc.Ctx.Err() != nil {
			ctxlog.Info(sc.Ctx, "autocommit: exiting", "err", sc.Ctx.Err())
			break
		}

		err := sc.findAndCommitVord()
		if err == nil {
			continue // on success, immediately try to commit another
		}

		if errors.Is(err, sql.ErrNoRows) {
			select {
			case <-sc.Ctx.Done():
			case <-time.After(time.Second * time.Duration(sc.Config.Vord.AutocommitPeriod.V)):
			}
			continue
		}

		ctxlog.Error2(sc.Ctx, "autocommit: error", err)
		time.Sleep(time.Second) // to avoid a busy loop when e.g. the database drops out
	}
}
