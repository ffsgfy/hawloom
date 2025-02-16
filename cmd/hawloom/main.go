package main

import (
	"context"
	"flag"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/handlers"
	"github.com/ffsgfy/hawloom/internal/config"
	"github.com/ffsgfy/hawloom/internal/utils"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func main() {
	ctxlog.SetDefault(ctxlog.New(os.Stdout, ctxlog.INFO))
	ctx := context.Background()

	configPath := flag.String("c", "", "config file path")
	flag.Parse()

	config, err := config.Load(*configPath)
	if err != nil {
		ctxlog.Error2(ctx, "failed to load config", err)
		panic(err)
	}

	logger := ctxlog.New(os.Stdout, ctxlog.Level(config.Misc.LogLevel.V))
	ctxlog.SetDefault(logger)
	ctx = ctxlog.WithLogger(ctx, logger)

	state, err := api.NewState(ctx, config)
	if err != nil {
		ctxlog.Error2(ctx, "failed to initialize state", err)
		panic(err)
	}

	if err := state.Ctx(ctx).LoadAuthKeys(true); err != nil {
		ctxlog.Error2(ctx, "failed to load keys", err)
		panic(err)
	}

	autocommitCtx, autocommitCancel := context.WithCancel(ctx)
	state.TasksCancel = append(state.TasksCancel, autocommitCancel)
	go state.Ctx(autocommitCtx).RunAutocommit()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.OFF)

	handlers.AddHandlers(e, state)

	utils.RunEcho(ctx, e, config.HTTP.BindAddress.V, config.HTTP.BindPort.V)

	ctxlog.Info(ctx, "waiting for background tasks")
	for _, cancel := range state.TasksCancel {
		cancel()
	}
	state.TasksWG.Wait()
}
