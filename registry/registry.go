package registry

import (
	"context"
	"google.golang.org/grpc"
)

type Registry interface {
	Register(context.Context, grpc.ServiceInfo)
	DeRegister()
}
