package configx

//go:generate enumer -type HandlerType -json -yaml -text -values -output logger_string.go
type HandlerType uint

const (
	HandlerTypeText HandlerType = iota
	HandlerTypeJson
)

type Logger struct {
	Level  string      `json:"level,omitempty" yaml:"level,omitempty" toml:"level,omitempty"`
	Prefix string      `json:"prefix,omitempty" yaml:"prefix,omitempty" toml:"prefix,omitempty"`
	Format HandlerType `json:"format,omitempty" yaml:"format,omitempty" toml:"format,omitempty"`
}

type standardLoggerConfig struct{}

func (f *standardLoggerConfig) Fetch() (Logger, error) {
	return *standard.Logger, nil
}
