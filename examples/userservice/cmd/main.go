package main

import (
	"context"
	"fmt"

	"github.com/charliego3/pallas"
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
			httpx.WithMiddleware(httpx.Middleware(func(next httpx.Handler) httpx.Handler {
				fmt.Println("into middleware 1....")
				return next
			})),
			httpx.WithMiddleware(httpx.RecoverMiddleware),
			httpx.WithMiddleware(httpx.Middleware(func(next httpx.Handler) httpx.Handler {
				return httpx.HandlerFunc(func(c *httpx.Context) error {
					fmt.Println(c.URL.Path)
					return next.Serve(c)
				})
			})),
		),
	}
}
