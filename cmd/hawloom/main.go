package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/handlers"
)

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUri := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbUser,
	)

	s, err := api.NewState(context.Background(), dbUri)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	handlers.AddHandlers(e, &s)

	run(e)
}

func run(e *echo.Echo) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	port := os.Getenv("HAWLOOM_PORT")
	if len(port) == 0 {
		e.Logger.Fatal(errors.New("server port not set"))
		cancel()
	}

	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(fmt.Errorf("failed to start server: %w", err))
			cancel()
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(fmt.Errorf("failed to shut down server: %w", err))
		e.Close()
	}
}
