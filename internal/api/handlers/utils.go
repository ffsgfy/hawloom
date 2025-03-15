package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	HXRedirect = "HX-Redirect"
	HXRefresh  = "HX-Refresh"
	HXRequest  = "HX-Request"
	// HXTrigger  = "HX-Trigger"
)

func isHTMX(c echo.Context) bool {
	return c.Request().Header.Get(HXRequest) != ""
}

func handleRedirect(c echo.Context, url string) error {
	if isHTMX(c) {
		c.Response().Header().Set(HXRedirect, url)
		return c.NoContent(http.StatusOK)
	}
	return c.Redirect(http.StatusSeeOther, url)
}

func handleRefresh(c echo.Context) error {
	if isHTMX(c) {
		c.Response().Header().Set(HXRefresh, "true")
		return c.NoContent(http.StatusOK)
	}
	if ref := c.Request().Referer(); ref != "" {
		return handleRedirect(c, ref)
	}
	return handleRedirect(c, "/")
}

func handleFormError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		var herr *echo.HTTPError
		if errors.As(err, &herr) {
			if herr.Code != http.StatusInternalServerError {
				if msg, ok := herr.Message.(string); ok {
					return errors.Join(err, c.HTML(http.StatusOK, "Error: "+msg))
				}
			}
		}
		return errors.Join(err, c.HTML(http.StatusOK, "Internal error"))
	}
}
