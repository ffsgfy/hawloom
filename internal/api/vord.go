package api

import (
	// "database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils"
)

type VordFlags int32

const (
	VordFlagTie        VordFlags = 1 << 1 // more than one ver with max votes
	VordFlagNoMajority VordFlags = 1 << 0 // majority rule not satisfied
)

func (sc *StateCtx) CreateVord(docID uuid.UUID, vordDuration time.Duration) error {
	res, err := sc.Queries.CreateVord(sc.Ctx, &db.CreateVordParams{
		Doc: docID,
		Duration: vordDuration,
	})
	if err != nil {
		return err
	}
	if res == 0 {
		return ErrVordExists
	}
	return nil
}

type CommitStatus int32

const (
	CommitStatusOK         CommitStatus = 0
	CommitStatusError      CommitStatus = 1 // unexpected error
	CommitStatusNoVotes    CommitStatus = 2 // no votes at all
	CommitStatusTie        CommitStatus = 3 // more than one ver with max votes
	CommitStatusNoMajority CommitStatus = 4 // majority rule not satisfied
)

// Assumes vord is already locked
func (sc *StateCtx) commitVord(vord *db.Vord, majorityRule bool) (CommitStatus, error) {
	if vord.Num != -1 {
		return CommitStatusError, ErrCommitPastVord
	}

	verVotes, err := sc.Queries.UpdateAllVerVotes(sc.Ctx, &db.UpdateAllVerVotesParams{
		Doc:     vord.Doc,
		VordNum: vord.Num,
	})
	if err != nil {
		return CommitStatusError, err
	}
	if len(verVotes) == 0 {
		return CommitStatusNoVotes, nil
	}

	var maxVotes int32
	var maxVotesTie bool
	for _, row := range verVotes {
		if row.Votes > maxVotes {
			maxVotes = row.Votes
			maxVotesTie = false
		} else if row.Votes == maxVotes {
			maxVotesTie = true
		}
	}

	if maxVotesTie {
		return CommitStatusTie, nil
	}

	if majorityRule {
		voters, err := sc.Queries.CountVoters(sc.Ctx, &db.CountVotersParams{
			Doc:     vord.Doc,
			VordNum: vord.Num,
		})
		if err != nil {
			return CommitStatusError, err
		}

		if int64(maxVotes) * 2 <= voters {
			return CommitStatusNoMajority, nil
		}
	}

	if err = sc.Queries.CommitVord(sc.Ctx, &db.CommitVordParams{
		Doc:   vord.Doc,
		Flags: 0,
	}); err != nil {
		return CommitStatusError, err
	}

	return CommitStatusOK, nil
}
