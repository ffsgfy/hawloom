package main

import (
	"context"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/api/handlers"
	"github.com/ffsgfy/hawloom/internal/utils"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func main() {
	ctxlog.SetDefault(ctxlog.New(os.Stdout, ctxlog.INFO))
	ctx := context.Background()

	portStr := os.Getenv("HAWLOOM_PORT")
	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		ctxlog.Error2(ctx, "invalid server port specified", err, "port", portStr)
		panic(err)
	}

	state, err := api.NewState(ctx, utils.MakePostgresURIFromEnv())
	if err != nil {
		ctxlog.Error2(ctx, "failed to initialize state", err)
		panic(err)
	}

	if err := state.Ctx(ctx).LoadKeys(true); err != nil {
		ctxlog.Error2(ctx, "failed to load keys", err)
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetLevel(log.OFF)

	handlers.AddHandlers(e, state)

	utils.RunEcho(ctx, e, uint16(port))
}
