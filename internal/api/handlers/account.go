package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
)

var (
	errNoAccountIDOrName    = echo.NewHTTPError(http.StatusBadRequest, "either account id or name required")
	errBothAccountIDAndName = echo.NewHTTPError(http.StatusBadRequest, "both account id and name not supported")
	errAccountNotFound      = echo.NewHTTPError(http.StatusNotFound, "account not found")
	errAccountWasNil        = echo.NewHTTPError(http.StatusInternalServerError, "account was nil")
	errPasswordTooLong      = echo.NewHTTPError(http.StatusBadRequest, "password too long")
	errAccountNameTaken     = echo.NewHTTPError(http.StatusConflict, "account name already taken")
)

type accountGetRequest struct {
	ID   *int32  `query:"id"`
	Name *string `query:"name"`
}

type accountGetResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func HandleAccountGet(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := accountGetRequest{}
		if err := c.Bind(&req); err != nil {
			return onBindError(err)
		}

		if req.ID == nil && req.Name == nil {
			return errNoAccountIDOrName
		}
		if req.ID != nil && req.Name != nil {
			return errBothAccountIDAndName
		}

		var err error
		var account *db.Account

		ctx := c.Request().Context()
		if req.ID != nil {
			account, err = s.Queries.FindAccountByID(ctx, *req.ID)
		}
		if req.Name != nil {
			account, err = s.Queries.FindAccountByName(ctx, *req.Name)
		}

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errAccountNotFound
			}
			return err
		}
		if account == nil {
			return errAccountWasNil
		}

		return c.JSON(http.StatusOK, &accountGetResponse{
			ID:        account.ID,
			Name:      account.Name,
			CreatedAt: account.CreatedAt.Time,
		})
	}
}

type accountPostRequest struct {
	Name     string `query:"name" form:"name"`
	Password string `query:"password" form:"password"`
}

type accountPostResponse struct {
	ID int32 `json:"id"`
}

func HandleAccountPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := accountPostRequest{}
		if err := c.Bind(&req); err != nil {
			return onBindError(err)
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				return errPasswordTooLong
			}
			return err
		}

		id, err := s.Queries.CreateAccount(c.Request().Context(), &db.CreateAccountParams{
			Name:         req.Name,
			PasswordHash: passwordHash,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errAccountNameTaken
			}
			return err
		}

		return c.JSON(http.StatusOK, &accountPostResponse{
			ID: id,
		})
	}
}
