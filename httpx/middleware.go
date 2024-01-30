package httpx

import "fmt"

type Middleware func(next Handler) Handler

func RecoverMiddleware(next Handler) Handler {
	return HandlerFunc(func(ctx *Context) (err error) {
		defer func() {
			if msg := recover(); msg != nil {
				err = fmt.Errorf("%v", msg)
			}
		}()

		return next.Serve(ctx)
	})
}
