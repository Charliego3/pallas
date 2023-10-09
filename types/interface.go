package types

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Server interface {
	Start(context.Context) error
	Run(context.Context) error
	Shutdown() error
	Logger() *slog.Logger
	Address() net.Addr
	RegisterService(svs ...Service)
}

type Service interface {
	Desc() *grpc.ServiceDesc
}
