package middleware

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

type reqIDKeyType struct{}

var reqIDKey = reqIDKeyType{}

func SetupContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqID := fmt.Sprintf("%016x", rand.Uint64())
		req := c.Request()
		ctx := context.WithValue(req.Context(), reqIDKey, reqID)
		c.SetRequest(req.WithContext(ctxlog.With(ctx, "req_id", reqID)))
		return next(c)
	}
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(reqIDKey).(string); ok {
		return reqID
	}
	return ""
}

func LogAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		req := c.Request()
		res := c.Response()
		start := time.Now()

		res.After(func() {
			duration := time.Since(start)
			durationMilli := float64(duration.Microseconds()) / 1000.0

			// TODO: also log user ip, identity
			tags := make([]any, 0, 10)
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

			ctxlog.Info(req.Context(), "HTTP "+req.Method, tags...)
		})

		err = next(c)
		if res.Committed && err != nil {
			ctxlog.Error2(req.Context(), "post-response error", err)
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
			"req_id":  GetRequestID(c.Request().Context()),
		}

		return herr
	}
}
