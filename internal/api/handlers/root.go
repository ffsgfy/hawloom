package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/middleware"
)

var (
	errInvalidInput = echo.NewHTTPError(http.StatusBadRequest, "invalid input")
)

func onBindError(err error) error {
	if herr, ok := err.(*echo.HTTPError); ok {
		err = herr.Internal
	}
	return errInvalidInput.WithInternal(err)
}

func AddHandlers(e *echo.Echo, s *api.State) {
	e.Use(middleware.SetRequestContext)
	e.GET("/healthcheck", HandleHealthcheck(s))

	g := e.Group("", middleware.LogAccess)
	g.GET("/account", HandleAccountGet(s))
	g.POST("/account", HandleAccountPost(s))
}
