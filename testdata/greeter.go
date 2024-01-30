package testdata

import (
	"context"

	"github.com/charliego3/pallas/types"
)

type Greeter struct {
	UnimplementedGreeterServer
}

func (g *Greeter) SayHello(_ context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: "reply with " + req.Name,
	}, nil
}

func (g *Greeter) SayHelloStream(Greeter_SayHelloStreamServer) error {
	return nil
}

func (g *Greeter) Desc() types.ServiceDesc {
	return types.ServiceDesc{
		Grpc: Greeter_ServiceDesc,
		Http: Greeter_HttpServiceDesc,
	}
}
