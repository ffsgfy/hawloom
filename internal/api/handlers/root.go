package handlers

import (
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
)

func AddHandlers(e *echo.Echo, s *api.State) {
	e.GET("/healthcheck", HandleHealthcheck(s))

	e.GET("/account", HandleAccountGet(s))
	e.POST("/account", HandleAccountPost(s))
}
