package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/ui"
)

func HandleMain(s *api.State) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := ui.Render(c.Request().Context(), ui.MainPage())
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, content)
	}
}
