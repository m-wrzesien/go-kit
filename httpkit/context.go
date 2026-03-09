package httpkit

import (
	"context"
)

type contextKey int

const (
	contextRemoteAddr contextKey = iota
	contextScheme
)

func GetRemoteAddr(ctx context.Context) string {
	return ctx.Value(contextRemoteAddr).(string)
}

func GetScheme(ctx context.Context) string {
	return ctx.Value(contextScheme).(string)
}
