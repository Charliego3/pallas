package types

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type ServiceRegister interface {
	RegisterService(svs ...Service)
}

type Server interface {
	ServiceRegister
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type BaseServer struct {
	Logger   *slog.Logger
	Listener net.Listener
}

type Dispatcher interface {
	Dispatch()
}

type Service interface {
	Desc() grpc.ServiceDesc
}
