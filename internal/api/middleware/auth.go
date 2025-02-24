package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
)

func ManageAuth(s *api.State) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authCookie, err := c.Cookie(s.Config.Auth.Cookie.V)
			if err == nil && len(authCookie.Value) > 0 {
				req := c.Request()
				authState := s.CreateAuthState(authCookie.Value)
				c.SetRequest(req.WithContext(api.WithAuthState(req.Context(), authState)))
			}

			// TODO: auto-renew tokens that will soon expire
			return next(c)
		}
	}
}
