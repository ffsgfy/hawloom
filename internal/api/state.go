package api

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffsgfy/hawloom/internal/db"
)

type State struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries

	KeysLock sync.RWMutex
	Keys     map[int32]*AuthKey
	KeyInUse int32
}

func NewState(ctx context.Context, dbURI string) (*State, error) {
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, err
	}

	return &State{
		Pool:    pool,
		Queries: db.New(pool),
		Keys:    map[int32]*AuthKey{},
	}, nil
}

type StateCtx struct {
	*State
	Ctx context.Context
}

func (s *State) Ctx(ctx context.Context) *StateCtx {
	return &StateCtx{State: s, Ctx: ctx}
}
