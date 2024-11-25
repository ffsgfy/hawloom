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

	g := e.Group("/", middleware.SetupContext, middleware.WrapErrors, middleware.LogAccess)
	g.GET("account", HandleAccountGet(s))
	g.GET("account/chkauth", HandleAccountGet(s), middleware.ManageAuth(s)) // TODO: remove this
	g.POST("account", HandleAccountPost(s))
	g.POST("account/login", HandleAccountLoginPost(s))
}
