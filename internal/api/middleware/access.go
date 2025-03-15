package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func SetupContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := api.CreateReqID()
		req := c.Request()
		c.SetRequest(
			req.WithContext(
				api.WithReqID(
					ctxlog.With(req.Context(), "req_id", reqID),
					reqID,
				),
			),
		)
		return next(c)
	}
}

func LogAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error

		start := time.Now()
		logfn := func() {
			duration := time.Since(start)
			durationMilli := float64(duration.Microseconds()) / 1000.0

			req := c.Request()
			res := c.Response()
			ctx := req.Context()

			// TODO: also log user ip
			tags := make([]any, 0, 16)
			tags = append(
				tags,
				"req_path", c.Path(),
				"res_status", res.Status,
				"res_size", res.Size,
				"dt", durationMilli,
			)

			if err != nil {
				tags = append(tags, "err", err)
			}

			authState := api.GetAuthState(ctx)
			if authState != nil {
				if authState.Token != nil {
					tags = append(tags, "auth_token", authState.Token)
				}
				if authState.Error != nil {
					tags = append(tags, "auth_err", authState.Error)
				}
			}

			ctxlog.Info(ctx, "HTTP "+req.Method, tags...)
		}

		err = next(c)
		if c.Response().Committed {
			logfn()
		} else {
			c.Response().After(logfn)
		}

		return err
	}
}

func WrapErrors(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		herr, ok := err.(*echo.HTTPError)
		if ok && herr.Code == http.StatusInternalServerError {
			// Replace all internal errors with the generic api.ErrInternal
			ok = false
		}
		if ok {
			if _, ok := herr.Internal.(*echo.HTTPError); ok {
				// Wrap internal HTTP error (see (echo.Echo).DefaultHTTPErrorHandler)
				herr = herr.WithInternal(errors.Join(herr.Internal))
			} else {
				// Just do a shallow copy
				herr = herr.WithInternal(herr.Internal)
			}
		} else {
			herr = api.ErrInternal.WithInternal(err)
		}

		herr.Message = map[string]any{
			"message": herr.Message,
			"req_id":  api.GetReqID(c.Request().Context()),
		}

		return herr
	}
}
