package handlers

import (
	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/middlewares"
)

func AddHandlers(e *echo.Echo, s *api.State) {
	e.Use(middlewares.SetRequestID)
	e.GET("/healthcheck", HandleHealthcheck(s))

	g := e.Group("", middlewares.LogAccess)
	g.GET("/account", HandleAccountGet(s))
	g.POST("/account", HandleAccountPost(s))
}
