package api

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffsgfy/hawloom/internal/db"
)

type State struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
	Auth    *Auth
}

func NewState(ctx context.Context, dbURI string) (*State, error) {
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, err
	}

	return &State{
		Pool:    pool,
		Queries: db.New(pool),
		Auth:    NewAuth(),
	}, nil
}

type StateCtx struct {
	*State
	Ctx context.Context
}

func (s *State) Ctx(ctx context.Context) *StateCtx {
	return &StateCtx{State: s, Ctx: ctx}
}

func (sc *StateCtx) TxWith(options pgx.TxOptions, fn func(*StateCtx) error) error {
	tx, err := sc.Pool.BeginTx(sc.Ctx, options)
	if err != nil {
		return err
	}
	defer tx.Rollback(sc.Ctx)

	if err = fn(&StateCtx{
		State: &State{
			Pool:    sc.Pool,
			Queries: sc.Queries.WithTx(tx),
			Auth:    sc.Auth,
		},
		Ctx: sc.Ctx,
	}); err != nil {
		return err
	}

	return tx.Commit(sc.Ctx)
}

func (sc *StateCtx) Tx(fn func(*StateCtx) error) error {
	return sc.TxWith(pgx.TxOptions{}, fn)
}
