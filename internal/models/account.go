package models

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Account struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	CreatedAt    time.Time `db:"created_at"`
	PasswordHash []byte    `db:"password_hash"`
}

var (
	AccountNameTakenError = errors.New("account name already taken")
	AccountNotFoundError  = errors.New("account not found")
)

func CreateAccount(db *sqlx.DB, name string, passwordHash []byte) (int, error) {
	account := Account{
		Name:         name,
		PasswordHash: passwordHash,
	}

	rows, err := db.NamedQuery(
		"INSERT INTO account (name, password_hash) "+
			"VALUES (:name, :password_hash) "+
			"ON CONFLICT DO NOTHING "+
			"RETURNING id",
		account,
	)
	if err != nil {
		return 0, err
	}
	if !rows.Next() {
		return 0, AccountNameTakenError
	}
	defer rows.Close()

	var id int
	if err = rows.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func FindAccountByID(db *sqlx.DB, id int) (*Account, error) {
	accounts := make([]Account, 0, 1)

	err := db.Select(&accounts, "SELECT * FROM account WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, AccountNotFoundError
	}
	return &accounts[0], nil
}

func FindAccountByName(db *sqlx.DB, name string) (*Account, error) {
	accounts := make([]Account, 0, 1)

	err := db.Select(&accounts, "SELECT * FROM account WHERE name=$1", name)
	if err != nil {
		return nil, err
	}

	if len(accounts) == 0 {
		return nil, AccountNotFoundError
	}
	return &accounts[0], nil
}
