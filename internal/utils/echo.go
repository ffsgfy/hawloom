package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func RunEcho(ctx context.Context, e *echo.Echo, port uint16) {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	ctxlog.Info(ctx, "server starting", "port", port)

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			ctxlog.Error2(ctx, "fatal server error", err)
			cancel()
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctxlog.Info(ctx, "server shutting down")

	if err := e.Shutdown(ctx); err != nil {
		ctxlog.Error2(ctx, "failed to shut down server", err)
		e.Close()
	}
}