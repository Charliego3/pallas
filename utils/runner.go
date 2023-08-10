package utils

type SafeRunner struct {
	err error
}

func NewSafeRunner() *SafeRunner {
	return &SafeRunner{}
}

func (r *SafeRunner) Run(fn func() error) {
	if r.err != nil {
		return
	}

	r.err = fn()
}

func (r *SafeRunner) Err() error {
	return r.err
}
