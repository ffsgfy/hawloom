package handlers

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/models"
)

type accountGetRequest struct {
	ID   *int    `query:"id"`
	Name *string `query:"name"`
}

type accountGetResponse struct {
	ID        int       `json:"id"`
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
		var account *models.Account

		if req.ID != nil {
			account, err = models.FindAccountByID(s.DB, *req.ID)
		}
		if req.Name != nil {
			account, err = models.FindAccountByName(s.DB, *req.Name)
		}

		if err != nil {
			if errors.Is(err, models.AccountNotFoundError) {
				return c.String(http.StatusNotFound, err.Error())
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if account == nil {
			return c.String(http.StatusInternalServerError, "account was nil")
		}

		return c.JSON(http.StatusOK, &accountGetResponse{
			ID:        account.ID,
			Name:      account.Name,
			CreatedAt: account.CreatedAt,
		})
	}
}

type accountPostRequest struct {
	Name     string `query:"name" form:"name"`
	Password string `query:"password" form:"password"`
}

type accountPostResponse struct {
	ID int `json:"id"`
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

		id, err := models.CreateAccount(s.DB, req.Name, passwordHash)
		if err != nil {
			if errors.Is(err, models.AccountNameTakenError) {
				return c.String(http.StatusConflict, err.Error())
			}
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &accountPostResponse{
			ID: id,
		})
	}
}
