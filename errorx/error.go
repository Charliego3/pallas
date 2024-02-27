package errorx

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	*status.Status
	casue error
}

func (e *Error) Error() string {
	return fmt.Sprintf("detail: %v", e.casue)
}

func New(code int, message string) *Error {
	return &Error{
		Status: status.New(codes.Canceled, message),
	}
}
