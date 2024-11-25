package main

import (
	"context"
	"os"

	"github.com/ffsgfy/hawloom/internal/api"
	"github.com/ffsgfy/hawloom/internal/utils"
	"github.com/ffsgfy/hawloom/internal/utils/ctxlog"
)

func main() {
	ctxlog.SetDefault(ctxlog.New(os.Stdout, ctxlog.INFO))
	ctx := context.Background()

	state, err := api.NewState(ctx, utils.MakePostgresURIFromEnv())
	if err != nil {
		ctxlog.Error2(ctx, "failed to initialize state", err)
		panic(err)
	}

	if err := state.Ctx(ctx).LoadKeys(false); err != nil {
		ctxlog.Error2(ctx, "failed to load keys", err)
		panic(err)
	}

	key, err := state.Ctx(ctx).CreateKey()
	if err != nil {
		ctxlog.Error2(ctx, "failed to create new key", err)
		panic(err)
	}

	ctxlog.Info(ctx, "created new key", "key_id", key.ID)
}
