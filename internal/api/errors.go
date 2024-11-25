package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// With the exception of ErrInternal, only add internal errors to this list
// if they are reused; for one-off internal errors use errors.New() or fmt.Errorf()
var (
	ErrInvalidInput = echo.NewHTTPError(http.StatusBadRequest, "invalid input")
	ErrInternal     = echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	ErrUnauthorized = echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")

	ErrNoAccountIDOrName       = echo.NewHTTPError(http.StatusBadRequest, "either account id or name required")
	ErrBothAccountIDAndName    = echo.NewHTTPError(http.StatusBadRequest, "both account id and name not supported")
	ErrAccountNotFound         = echo.NewHTTPError(http.StatusNotFound, "account not found")
	ErrAccountNameTooShort     = echo.NewHTTPError(http.StatusBadRequest, "account name too short")
	ErrAccountPasswordTooShort = echo.NewHTTPError(http.StatusBadRequest, "account password too short")
	ErrPasswordTooLong         = echo.NewHTTPError(http.StatusBadRequest, "password too long")
	ErrAccountNameTaken        = echo.NewHTTPError(http.StatusConflict, "account name already taken")

	ErrMalformedToken = echo.NewHTTPError(http.StatusUnauthorized, "malformed token")
	ErrNoTokenKey     = echo.NewHTTPError(http.StatusUnauthorized, "non-existent token key")
	ErrWrongTokenHash = echo.NewHTTPError(http.StatusUnauthorized, "wrong token hash")
	ErrExpiredToken   = echo.NewHTTPError(http.StatusUnauthorized, "expired token")
)

func OnBindError(err error) error {
	if herr, ok := err.(*echo.HTTPError); ok {
		err = herr.Internal
	}
	return ErrInvalidInput.WithInternal(err)
}
