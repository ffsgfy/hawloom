package api

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
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

	return account, nil
}

func (s *State) checkAccountNamePasswordLengths(name, password string) error {
	if len(name) < s.Config.Account.NameMinLength.V {
		return ErrAccountNameTooShort
	}
	if len(name) > s.Config.Account.NameMaxLength.V {
		return ErrAccountNameTooLong
	}
	if len(password) < s.Config.Account.PasswordMinLength.V {
		return ErrAccountPasswordTooShort
	}
	if len(password) > s.Config.Account.PasswordMaxLength.V {
		return ErrAccountPasswordTooLong
	}
	return nil
}

func (sc *StateCtx) CreateAccount(name, password string) (*db.Account, error) {
	if err := sc.checkAccountNamePasswordLengths(name, password); err != nil {
		return nil, err
	}

	// Check if name exists before hashing the password
	if exists, err := sc.Queries.CheckAccountName(sc.Ctx, name); err != nil {
		return nil, err
	} else if exists != 0 {
		return nil, ErrAccountNameTaken
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return nil, ErrAccountPasswordTooLong
		}
		return nil, err
	}

	account, err := sc.Queries.CreateAccount(sc.Ctx, name, passwordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNameTaken
		}
		return nil, err
	}

	ctxlog.Info(
		sc.Ctx, "account created",
		"account_id", account.ID,
		"account_name", account.Name,
	)

	return account, nil
}

func (sc *StateCtx) CheckPassword(name, password string) (*db.Account, error) {
	if err := sc.checkAccountNamePasswordLengths(name, password); err != nil {
		return nil, err
	}

	account, err := sc.FindAccount(nil, &name)
	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword(account.PasswordHash, []byte(password)); err != nil {
		return nil, ErrUnauthorized.WithInternal(err)
	}

	ctxlog.Info(
		sc.Ctx, "account password matched",
		"account_id", account.ID,
		"account_name", account.Name,
	)

	return account, nil
}
