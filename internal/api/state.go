package api

import (
	"github.com/jmoiron/sqlx"
)

type State struct {
	Db *sqlx.DB
}

func NewState(dbDriver, dbUri string) (State, error) {
	db, err := sqlx.Connect(dbDriver, dbUri)
	if err != nil {
		return State{}, err
	}

	s := State{
		Db: db,
	}

	return s, nil
}
