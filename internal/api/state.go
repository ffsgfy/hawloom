package api

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ffsgfy/hawloom/internal/db"
)

type State struct {
	Pool    *pgxpool.Pool
	Queries *db.Queries
}

func NewState(ctx context.Context, dbUri string) (State, error) {
	pool, err := pgxpool.New(ctx, dbUri)
	if err != nil {
		return State{}, err
	}

	s := State{
		Pool:    pool,
		Queries: db.New(pool),
	}

	return s, nil
}
