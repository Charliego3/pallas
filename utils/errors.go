package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorArgument(err error) error {
	if e, ok := err.(interface{ Field() string }); ok {
		return status.Errorf(codes.InvalidArgument, "%s verification failed", e.Field())
	}
	return err
}

func ErrorArgumentf(format string, v ...any) error {
	if len(v) == 0 {
		return status.Error(codes.InvalidArgument, format)
	}
	return status.Errorf(codes.InvalidArgument, format, v...)
}
