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

	baseGroup := e.Group("", middleware.SetupContext, middleware.WrapErrors, middleware.LogAccess)
	authGroup := baseGroup.Group("", middleware.ManageAuth(s))

	authGroup.GET("/", HandleMain(s))
	authGroup.GET("/user/:name", HandleUser(s))

	authGroup.GET("/auth/login", HandleLogin(s))
	baseGroup.POST("/auth/login", HandleLoginPost(s), handleFormError)
	authGroup.GET("/auth/register", HandleRegister(s))
	baseGroup.POST("/auth/register", HandleRegisterPost(s), handleFormError)
	authGroup.GET("/auth/logout", HandleLogout(s))

	authGroup.GET("/doc/new", HandleNewDoc(s))
	authGroup.POST("/doc/new", HandleNewDocPost(s), handleFormError)
	authGroup.GET("/doc/:doc", HandleDoc(s))
	authGroup.GET("/doc/:doc/:vord", HandleDoc(s))
	authGroup.GET("/ver/list", HandleVerList(s))
	authGroup.GET("/ver/:ver", HandleVer(s))
	authGroup.POST("/ver/:ver/vote", HandleVerVoteUnvote(s, true))
	authGroup.POST("/ver/:ver/unvote", HandleVerVoteUnvote(s, false))
	authGroup.DELETE("/ver/:ver", HandleVerDelete(s))
}
