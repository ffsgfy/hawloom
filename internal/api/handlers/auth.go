package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/db"
	"github.com/ffsgfy/hawloom/internal/ui"
)

func handleAuthError(c echo.Context, err error) error {
	var herr *echo.HTTPError
	if errors.As(err, &herr) {
		if herr.Code != http.StatusInternalServerError {
			if msg, ok := herr.Message.(string); ok {
				return c.HTML(http.StatusOK, "Error: "+msg)
			}
		}
	}
	return err
}

func handleAuthSuccess(c echo.Context, sc *api.StateCtx, account *db.Account) error {
	key := sc.Auth.KeyInUse.Load()
	token := api.CreateAuthToken(key, account.Name, account.ID, sc.Config.Auth.TokenTTL.V)
	cookie, err := api.CreateAuthCookie(key, token, sc.Config.Auth.Cookie.V)
	if err != nil {
		return err
	}

	c.SetCookie(cookie)
	c.Response().Header().Set(HXRedirect, "/")
	return c.NoContent(http.StatusOK)
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
			return handleAuthError(c, err)
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
			return handleAuthError(c, echo.NewHTTPError(http.StatusBadRequest, "passwords do not match"))
		}

		sc := s.Ctx(c.Request().Context())
		account, err := sc.CreateAccount(params.Name, params.Password)
		if err != nil {
			return handleAuthError(c, err)
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
		c.Response().Header().Set(HXRefresh, "true")
		return c.NoContent(http.StatusOK)
	}
}
