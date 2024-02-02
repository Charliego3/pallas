package main

import (
	"context"

	"github.com/charliego3/pallas"
	_ "github.com/charliego3/pallas/encoding/json"
	"github.com/charliego3/pallas/examples/userservice/internal/service"
	"github.com/charliego3/pallas/httpx"
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
		pallas.WithHttpOpts(
			httpx.WithMiddleware(httpx.RecoverMiddleware),
		),
	}
}
