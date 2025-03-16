package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
)

func ManageAuth(s *api.State) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cookie, err := c.Cookie(s.Config.Auth.Cookie.V); err == nil && len(cookie.Value) > 0 {
				req := c.Request()
				authState := s.CreateAuthState(cookie.Value)
				ctx := api.WithAuthState(req.Context(), authState)
				c.SetRequest(req.WithContext(ctx))

				if authState.Valid() && authState.Token.TTL(s.Config.Auth.TokenRenewTTL.V) <= 0 {
					if cookie, err := s.Ctx(ctx).CreateAuthCookie(
						authState.Token.AccountName, authState.Token.AccountID,
					); err == nil {
						c.SetCookie(cookie)
					}
				}
			}

			// TODO: auto-renew tokens that will soon expire
			return next(c)
		}
	}
}
