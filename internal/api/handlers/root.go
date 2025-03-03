package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/middleware"
)

const (
	HXRedirect = "HX-Redirect"
	HXRefresh  = "HX-Refresh"
)

func HandleNotFound(c echo.Context) error {
	return c.NoContent(http.StatusNotFound)
}

func AddHandlers(e *echo.Echo, s *api.State) {
	e.RouteNotFound("*", HandleNotFound)
	e.GET("/healthcheck", HandleHealthcheck(s))
	e.Static("/static", "static")

	baseGroup := e.Group("", middleware.SetupContext, middleware.WrapErrors, middleware.LogAccess)
	authGroup := baseGroup.Group("", middleware.ManageAuth(s))

	authGroup.GET("/", HandleMain(s))

	authGroup.GET("/auth/login", HandleLogin(s))
	baseGroup.POST("/auth/login", HandleLoginPost(s))
	authGroup.GET("/auth/register", HandleRegister(s))
	baseGroup.POST("/auth/register", HandleRegisterPost(s))
	authGroup.GET("/auth/logout", HandleLogout(s))
}
