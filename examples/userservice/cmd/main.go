package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/charliego3/pallas"
	"github.com/charliego3/pallas/examples/userservice/internal/service"
	"github.com/charliego3/pallas/middleware"
	"github.com/charliego3/pallas/middleware/logging"
	"github.com/charliego3/pallas/utility"
)

func main() {
	app := CreateApplication()
	app.RegisterService(
		new(service.Greeter),
		new(service.User),
	)
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}

func appOpts() []utility.Option[pallas.Application] {
	return []utility.Option[pallas.Application]{
		pallas.WithTCPAddr(":8888"),
		pallas.WithMiddleware(
			logging.Server(slog.Default()),
			func(next middleware.Handler) middleware.Handler {
				return func(ctx *middleware.Context) (any, error) {
					fmt.Println("middleware with app", ctx.Method, ctx.Path, ctx.Kind)
					ctx.ResHeader.Add("User-Server", "pallas")
					return next(ctx)
				}
			},
			func(next middleware.Handler) middleware.Handler {
				return func(ctx *middleware.Context) (any, error) {
					fmt.Println("middleware with app 2", ctx.ReqHeader)
					return next(ctx)
				}
			},
		),
	}
}
