package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
)

func handleAuthSuccess(c echo.Context, sc *api.StateCtx, account *db.Account) error {
	cookie, err := sc.CreateAuthCookie(account.Name, account.ID)
	if err != nil {
		return err
	}

	c.SetCookie(cookie)
	return handleRedirect(c, "/")
}

func HandleLogin(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := ui.Render(c.Request().Context(), ui.LoginPage())
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type loginParams struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

func HandleLoginPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params loginParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		sc := s.Ctx(c.Request().Context())
		account, err := sc.CheckPassword(params.Name, params.Password)
		if err != nil {
			return err
		}

		return handleAuthSuccess(c, sc, account)
	}
}

func HandleRegister(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := ui.Render(c.Request().Context(), ui.RegisterPage())
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}

type registerParams struct {
	Name       string `form:"name"`
	Password   string `form:"password"`
	PasswordRe string `form:"password-re"`
}

func HandleRegisterPost(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params registerParams
		if err := c.Bind(&params); err != nil {
			return err
		}

		if params.Password != params.PasswordRe {
			return echo.NewHTTPError(http.StatusBadRequest, "passwords do not match")
		}

		sc := s.Ctx(c.Request().Context())
		account, err := sc.CreateAccount(params.Name, params.Password)
		if err != nil {
			return err
		}

		return handleAuthSuccess(c, sc, account)
	}
}

func HandleLogout(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:    s.Config.Auth.Cookie.V,
			Path:    "/",
			Expires: time.Unix(0, 0),
		})
		return handleRefresh(c)
	}
}
