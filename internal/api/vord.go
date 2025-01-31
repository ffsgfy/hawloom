package api

import (
	// "database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/ffsgfy/hawloom/internal/db"
)

type VordFlags int32

const (
	VordFlagTie        VordFlags = 1 << 1 // more than one ver with max votes
	VordFlagNoMajority VordFlags = 1 << 0 // majority rule not satisfied
)

func (sc *StateCtx) CreateVord(docID uuid.UUID, vordDuration time.Duration) error {
	res, err := sc.Queries.CreateVord(sc.Ctx, docID, vordDuration)
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

	verVotes, err := sc.Queries.FindVersForCommit(sc.Ctx, vord.Doc)
	if err != nil {
		return CommitStatusError, err
	}
	if len(verVotes) == 0 || verVotes[0].Votes == 0 {
		return CommitStatusNoVotes, nil
	}
	if len(verVotes) > 1 && verVotes[0].Votes == verVotes[1].Votes {
		return CommitStatusTie, nil
	}

	if majorityRule {
		voters, err := sc.Queries.CountVoters(sc.Ctx, vord.Doc, vord.Num)
		if err != nil {
			return CommitStatusError, err
		}
		if int64(verVotes[0].Votes)*2 <= voters {
			return CommitStatusNoMajority, nil
		}
	}

	if err = sc.Queries.CommitVord(sc.Ctx, vord.Doc, 0); err != nil {
		return CommitStatusError, err
	}

	return CommitStatusOK, nil
}
