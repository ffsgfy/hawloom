package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/middleware"
)

func HandleNotFound(c echo.Context) error {
	return c.NoContent(http.StatusNotFound)
}

func AddHandlers(e *echo.Echo, s *api.State) {
	e.RouteNotFound("*", HandleNotFound)
	e.GET("/healthcheck", HandleHealthcheck(s))
	e.Static("/static", "static")

	g := e.Group("/", middleware.SetupContext, middleware.WrapErrors, middleware.LogAccess)
	g.GET("", HandleMain(s))
}
