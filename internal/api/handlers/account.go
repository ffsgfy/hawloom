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
			CreatedAt: account.CreatedAt.Time,
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

		key := sc.GetKeyInUse()
		if key == nil {
			return api.ErrNoKeyInUse
		}

		token := key.CreateToken(account.ID)
		tokenStr, err := key.EncodeToken(token)
		if err != nil {
			return err
		}

		c.SetCookie(&http.Cookie{
			Name:     api.AuthCookie,
			Value:    tokenStr,
			MaxAge:   api.TokenTTL,
			SameSite: http.SameSiteStrictMode,
		})

		return c.JSON(http.StatusOK, &accountPostResponse{
			ID: account.ID,
		})
	}
}
