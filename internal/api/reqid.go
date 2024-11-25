package api

import (
	"context"
	"fmt"
	"math/rand"
)

type reqIDKeyType struct{}

var reqIDKey = reqIDKeyType{}

func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

func GetReqID(ctx context.Context) string {
	if reqID, ok := ctx.Value(reqIDKey).(string); ok {
		return reqID
	}
	return ""
}

func CreateReqID() string {
	return fmt.Sprintf("%016x", rand.Uint64())
}
