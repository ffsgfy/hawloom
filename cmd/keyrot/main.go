package main

import (
	"context"
	"flag"
	"os"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/config"
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

	state, err := api.NewState(ctx, config)
	if err != nil {
		ctxlog.Error2(ctx, "failed to initialize state", err)
		panic(err)
	}

	if err := state.Ctx(ctx).LoadAuthKeys(false); err != nil {
		ctxlog.Error2(ctx, "failed to load keys", err)
		panic(err)
	}

	key, err := state.Ctx(ctx).CreateAuthKey()
	if err != nil {
		ctxlog.Error2(ctx, "failed to create new key", err)
		panic(err)
	}

	ctxlog.Info(ctx, "created new key", "key_id", key.ID)
}
