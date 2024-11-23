package api

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/ffsgfy/hawloom/internal/db"
)

func (sc *StateCtx) FindAccount(id *int32, name *string) (*db.Account, error) {
	if id == nil && name == nil {
		return nil, ErrNoAccountIDOrName
	}
	if id != nil && name != nil {
		return nil, ErrBothAccountIDAndName
	}

	var err error
	var account *db.Account

	if id != nil {
		account, err = sc.Queries.FindAccountByID(sc.Ctx, *id)
	}
	if name != nil {
		account, err = sc.Queries.FindAccountByName(sc.Ctx, *name)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountWasNil
	}

	return account, nil
}

func (sc *StateCtx) CreateAccount(name, password string) (int32, error) {
	if len(name) < 4 {
		return 0, ErrAccountNameTooShort
	}

	// Check if name exists before hashing the password
	exists, err := sc.Queries.CheckAccountName(sc.Ctx, name)
	if err != nil {
		return 0, err
	} else if exists != 0 {
		return 0, ErrAccountNameTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return 0, ErrPasswordTooLong
		}
		return 0, err
	}

	id, err := sc.Queries.CreateAccount(sc.Ctx, &db.CreateAccountParams{
		Name:         name,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrAccountNameTaken
		}
		return 0, err
	}

	return id, nil
}
