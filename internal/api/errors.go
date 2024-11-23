package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	ErrInvalidInput = echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	ErrInternal     = echo.NewHTTPError(http.StatusBadRequest, "internal error")

	ErrNoAccountIDOrName    = echo.NewHTTPError(http.StatusBadRequest, "either account id or name required")
	ErrBothAccountIDAndName = echo.NewHTTPError(http.StatusBadRequest, "both account id and name not supported")
	ErrAccountNotFound      = echo.NewHTTPError(http.StatusNotFound, "account not found")
	ErrAccountWasNil        = echo.NewHTTPError(http.StatusInternalServerError, "account was nil")
	ErrAccountNameTooShort  = echo.NewHTTPError(http.StatusBadRequest, "account name too short")
	ErrPasswordTooLong      = echo.NewHTTPError(http.StatusBadRequest, "password too long")
	ErrAccountNameTaken     = echo.NewHTTPError(http.StatusConflict, "account name already taken")
)

func OnBindError(err error) error {
	if herr, ok := err.(*echo.HTTPError); ok {
		err = herr.Internal
	}
	return ErrInvalidInput.WithInternal(err)
}
