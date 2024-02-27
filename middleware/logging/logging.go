package logging

import (
	"context"
	"github.com/charliego3/pallas/utility"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/charliego3/pallas/middleware"
)

func Server(logger *slog.Logger) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		return func(ctx *middleware.Context) (any, error) {
			startTime := time.Now()
			reply, err := next(ctx)
			attrs := []any{
				slog.String("kind", string(ctx.Kind)),
				slog.String("path", ctx.Path),
				slog.Any("req", ctx.Payload),
			}
			if utility.NonBlank(ctx.Method) {
				attrs = append(attrs, slog.String("method", ctx.Method))
			}
			if len(ctx.ResHeader) > 0 {
				attrs = append(attrs, slog.Any("reqHeader", ctx.ReqHeader))
			}
			if len(ctx.ResHeader) > 0 {
				attrs = append(attrs, slog.Any("resHeader", ctx.ResHeader))
			}
			level := slog.LevelInfo
			if err != nil {
				level = slog.LevelError
				attrs = append(attrs,
					slog.Any("err", err),
					slog.String("stack", string(debug.Stack())),
				)
			}
			attrs = append(attrs, slog.Duration("took", time.Since(startTime)))
			logger.Log(context.Background(), level, "request", attrs...)
			return reply, err
		}
	}
}
