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

func NewDefaultBaseServer() *BaseServer {
	s := new(BaseServer)
	s.Logger = slog.Default()
	return s
}

type Dispatcher interface {
	Dispatch()
}

type GrpcServiceDesc = grpc.ServiceDesc

type HttpMethodDesc struct {
	Method   string
	Template string
	Handler  func(Service) any
}

type HttpServiceDesc struct {
	ServiceName string
	HandlerType any
	Methods     []HttpMethodDesc
}

type ServiceDesc struct {
	Grpc GrpcServiceDesc
	Http HttpServiceDesc
}

type Service interface {
	Desc() ServiceDesc
}
