package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
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
			return api.OnBindError(err)
		}

		account, err := s.Ctx(c.Request().Context()).FindAccount(req.ID, req.Name)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &accountGetResponse{
			ID:        account.ID,
			Name:      account.Name,
			CreatedAt: account.CreatedAt,
		})
	}
}

type accountPostRequest struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

type accountPostResponse struct {
	ID int32 `json:"id"`
}

func HandleAccountPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := accountPostRequest{}
		if err := c.Bind(&req); err != nil {
			return api.OnBindError(err)
		}

		account, err := s.Ctx(c.Request().Context()).CreateAccount(req.Name, req.Password)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &accountPostResponse{
			ID: account.ID,
		})
	}
}

func HandleAccountLoginPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := accountPostRequest{}
		if err := c.Bind(&req); err != nil {
			return api.OnBindError(err)
		}

		sc := s.Ctx(c.Request().Context())
		account, err := sc.CheckPassword(req.Name, req.Password)
		if err != nil {
			return err
		}

		key := sc.Auth.KeyInUse.Load()
		token := api.CreateAuthToken(key, account.Name, account.ID, sc.Config.Auth.TokenTTL.V)
		cookie, err := api.CreateAuthCookie(key, token, sc.Config.Auth.Cookie.V)
		if err != nil {
			return err
		}

		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, &accountPostResponse{
			ID: account.ID,
		})
	}
}
