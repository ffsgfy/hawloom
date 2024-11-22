package middleware

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func SetRequestContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		reqid := fmt.Sprintf("%016x", rand.Uint64())
		req := c.Request()
		c.SetRequest(req.WithContext(ctxlog.With(req.Context(), "reqid", reqid)))
		return next(c)
	}
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
				"path", c.Path(),
				"status", res.Status,
				"dt", durationMilli,
				"size", res.Size,
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
