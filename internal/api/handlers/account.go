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
			return c.String(http.StatusBadRequest, err.Error())
		}

		if req.ID == nil && req.Name == nil {
			return c.String(http.StatusBadRequest, "either account id or name required")
		}
		if req.ID != nil && req.Name != nil {
			return c.String(http.StatusBadRequest, "both account id and name not supported")
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
				return c.String(http.StatusNotFound, "account not found")
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if account == nil {
			return c.String(http.StatusInternalServerError, "account was nil")
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
			return c.String(http.StatusBadRequest, err.Error())
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				return c.String(http.StatusBadRequest, "password too long")
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}

		id, err := s.Queries.CreateAccount(c.Request().Context(), &db.CreateAccountParams{
			Name:         req.Name,
			PasswordHash: passwordHash,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.String(http.StatusConflict, "account name already taken")
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &accountPostResponse{
			ID: id,
		})
	}
}
