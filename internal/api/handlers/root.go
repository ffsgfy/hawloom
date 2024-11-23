package handlers

import (
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/middleware"
)

func AddHandlers(e *echo.Echo, s *api.State) {
	e.Use(middleware.SetupContext)
	e.GET("/healthcheck", HandleHealthcheck(s))

	g := e.Group("", middleware.WrapErrors, middleware.LogAccess)
	g.GET("/account", HandleAccountGet(s))
	g.POST("/account", HandleAccountPost(s))
}
