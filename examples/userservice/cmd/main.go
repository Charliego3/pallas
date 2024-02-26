package main

import (
	"context"
	"fmt"

	"github.com/charliego3/pallas"
	"github.com/charliego3/pallas/examples/userservice/internal/service"
	"github.com/charliego3/pallas/middleware"
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
			func(ctx *middleware.Context) (any, error) {
				fmt.Println("middleware with app", ctx.Method, ctx.Path, ctx.Kind)
				ctx.ResHeader.Add("User-Server", "pallas")
				return nil, fmt.Errorf("error string")
			},
			func(ctx *middleware.Context) (any, error) {
				fmt.Println("middleware with app 2", ctx.ReqHeader)
				return nil, nil
			},
		),
	}
}
