package testdata

import (
	"context"
	"google.golang.org/grpc"
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

func (g *Greeter) Desc() grpc.ServiceDesc {
	return Greeter_ServiceDesc
}
