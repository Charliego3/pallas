package types

import (
	"context"
	"github.com/charliego3/mspp/utility"
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

type ServerInfo interface {
	Listener() net.Listener
}

type BaseServer struct {
	Network, Addr string
	Logger        *slog.Logger
	Listener      net.Listener
}

func (s *BaseServer) HasListener() bool {
	if s.Listener != nil {
		return true
	}

	return utility.NonBlanks(s.Network, s.Addr)
}

type Dispatcher interface {
	Dispatch()
}

type Service interface {
	Desc() *grpc.ServiceDesc
}
