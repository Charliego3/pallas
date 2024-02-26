package grpcx

import (
	"context"
	"fmt"

	"github.com/charliego3/pallas/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Printf("interceptor: %+v\n", req)
	context := middleware.NewGRPCContext(ctx, info.FullMethod, req)
	defer func() {
		fmt.Println(context.ResHeader)
		if len(context.ResHeader) == 0 {
			return
		}

		grpc.SetHeader(ctx, metadata.MD(context.ResHeader))
	}()
	for _, m := range s.middlewares {
		reply, err := m(context)
		if err != nil || reply != nil {
			return reply, err
		}
	}
	return handler(ctx, req)
}

func (s *Server) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	grpc.MethodFromServerStream(ss)
	metadata.FromIncomingContext(ss.Context())
	return handler(srv, ss)
}
