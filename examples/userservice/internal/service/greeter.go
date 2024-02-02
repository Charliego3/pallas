package service

import (
	"context"

	"github.com/charliego3/pallas/examples/protos"
)

type Greeter struct {
	pb.UnimplementedGreeterDescServer
}

func (g *Greeter) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "reply with " + req.Name,
	}, nil
}

func (g *Greeter) SayHelloStream(pb.Greeter_SayHelloStreamServer) error {
	return nil
}
