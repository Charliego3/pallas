package grpcx

import (
	"context"
	"github.com/charliego3/pallas/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	gctx := middleware.NewGRPCContext(ctx, info.FullMethod, req)
	m := middleware.Chain(s.middlewares...)
	return m(func(mctx *middleware.Context) (any, error) {
		reply, err := handler(mctx, req)
		_ = grpc.SetHeader(ctx, metadata.MD(gctx.ResHeader))
		return reply, err
	})(gctx)
}

func (s *Server) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	grpc.MethodFromServerStream(ss)
	metadata.FromIncomingContext(ss.Context())
	return handler(srv, ss)
}
