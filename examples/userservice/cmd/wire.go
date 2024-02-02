//go:build wireinject
// +build wireinject

package main

import (
	"github.com/charliego3/pallas"
	"github.com/google/wire"
)

func CreateApplication() *pallas.Application {
	wire.Build(pallas.NewApp, appOpts)
	return nil
}
