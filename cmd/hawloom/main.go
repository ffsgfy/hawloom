package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/handlers"
)

func main() {
	s, err := api.NewState("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/healthcheck", handlers.HandleHealthcheck(&s))

	run(e)
}

func run(e *echo.Echo) {
	sigintChan := make(chan os.Signal, 1)
	signal.Notify(sigintChan, os.Interrupt)

	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(fmt.Errorf("shutting down: %w", err))
		}
	}()

	<-sigintChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
		e.Close()
	}
}
