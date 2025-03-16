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
	ErrForbidden    = echo.NewHTTPError(http.StatusForbidden, "forbidden")

	ErrNoAccountIDOrName       = echo.NewHTTPError(http.StatusBadRequest, "either account id or name required")
	ErrBothAccountIDAndName    = echo.NewHTTPError(http.StatusBadRequest, "both account id and name not supported")
	ErrAccountNotFound         = echo.NewHTTPError(http.StatusNotFound, "account not found")
	ErrAccountNameTooShort     = echo.NewHTTPError(http.StatusBadRequest, "account name too short")
	ErrAccountNameTooLong      = echo.NewHTTPError(http.StatusBadRequest, "account name too long")
	ErrAccountPasswordTooShort = echo.NewHTTPError(http.StatusBadRequest, "account password too short")
	ErrAccountPasswordTooLong  = echo.NewHTTPError(http.StatusBadRequest, "account password too long")
	ErrAccountNameTaken        = echo.NewHTTPError(http.StatusConflict, "account name already taken")

	ErrMalformedToken = echo.NewHTTPError(http.StatusUnauthorized, "malformed token")
	ErrNoTokenKey     = echo.NewHTTPError(http.StatusUnauthorized, "non-existent token key")
	ErrWrongTokenHash = echo.NewHTTPError(http.StatusUnauthorized, "wrong token hash")
	ErrExpiredToken   = echo.NewHTTPError(http.StatusUnauthorized, "expired token")

	ErrDocTitleTooShort      = echo.NewHTTPError(http.StatusBadRequest, "document title too short")
	ErrDocTitleTooLong       = echo.NewHTTPError(http.StatusBadRequest, "document title too long")
	ErrRoundDurationTooSmall = echo.NewHTTPError(http.StatusBadRequest, "voting round duration too small")
	ErrDocNotFound           = echo.NewHTTPError(http.StatusNotFound, "document not found")
	ErrVerNotFound           = echo.NewHTTPError(http.StatusNotFound, "version not found")
	ErrDeletePastVer         = echo.NewHTTPError(http.StatusConflict, "cannot delete past version")

	ErrVotePastVer    = echo.NewHTTPError(http.StatusConflict, "cannot vote for past version")
	ErrVerVoteExists  = echo.NewHTTPError(http.StatusConflict, "already voted for this version")
	ErrDocVoteExists  = echo.NewHTTPError(http.StatusConflict, "already voted for another version")
	ErrDeletePastVote = echo.NewHTTPError(http.StatusConflict, "cannot delete past vote")
	ErrVoteNotFound   = echo.NewHTTPError(http.StatusNotFound, "vote not found")

	ErrNoVordExists   = echo.NewHTTPError(http.StatusConflict, "no active voting round currently exists")
	ErrVordExists     = echo.NewHTTPError(http.StatusConflict, "active voting round already exists")
	ErrCommitPastVord = echo.NewHTTPError(http.StatusConflict, "cannot commit past voting round")
)

func OnBindError(err error) error {
	if herr, ok := err.(*echo.HTTPError); ok {
		err = herr.Internal
	}
	return ErrInvalidInput.WithInternal(err)
}
