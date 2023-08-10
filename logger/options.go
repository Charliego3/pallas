package logger

import (
	"github.com/charliego3/mspp/opts"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"time"
)

// Replacer is called to rewrite each non-group attribute before it is logged.
// The attribute's value has been resolved (see [Value.Resolve]).
// If Replacer returns an Attr with Key == "", the attribute is discarded.
//
// The built-in attributes with keys "time", "level", "source", and "msg"
// are passed to this function, except that time is omitted
// if zero, and source is omitted if AddSource is false.
//
// The first argument is a list of currently open groups that contain the
// Attr. It must not be retained or modified. Replacer is never called
// for Group attributes, only their contents. For example, the attribute
// list
//
//	Int("a", 1), Group("g", Int("b", 2)), Int("c", 3)
//
// results in consecutive calls to Replacer with the following arguments:
//
//	nil, Int("a", 1)
//	[]string{"g"}, Int("b", 2)
//	nil, Int("c", 3)
//
// Replacer can be used to change the default keys of the built-in
// attributes, convert types (for example, to replace a `time.Time` with the
// integer seconds since the Unix epoch), sanitize personal information, or
// remove attributes from the output.
type Replacer func(groups []string, a slog.Attr) slog.Attr

type HandlerOptions struct {
	// timeFormat specify what's pattern to be formated
	// default using time.Kitchen
	//
	// eg: time.DateTime
	timeFormat string

	// w is output writer, default using os.Stderr
	w io.Writer

	// level is logger min Level, default is slog.LevelInfo
	level slog.Level

	// prefix output prefix in every record
	prefix string

	// replacer report to Replacer
	replacer Replacer
}

func getOptions(opts ...opts.Option[HandlerOptions]) *HandlerOptions {
	options := &HandlerOptions{
		timeFormat: time.Kitchen,
		w:          os.Stderr,
		level:      slog.LevelInfo,
	}
	for _, opt := range opts {
		opt.Apply(options)
	}
	return options
}

func WithTimeFormat(format string) opts.Option[HandlerOptions] {
	return opts.OptionFunc[HandlerOptions](func(cfg *HandlerOptions) {
		cfg.timeFormat = format
	})
}

func WithOutput(w io.Writer) opts.Option[HandlerOptions] {
	return opts.OptionFunc[HandlerOptions](func(cfg *HandlerOptions) {
		cfg.w = w
	})
}

func WithLevel(level slog.Level) opts.Option[HandlerOptions] {
	return opts.OptionFunc[HandlerOptions](func(cfg *HandlerOptions) {
		cfg.level = level
	})
}

func WithPrefix(prefix string) opts.Option[HandlerOptions] {
	return opts.OptionFunc[HandlerOptions](func(cfg *HandlerOptions) {
		cfg.prefix = prefix
	})
}

// WithReplacer please refer to Replacer
func WithReplacer(fn Replacer) opts.Option[HandlerOptions] {
	return opts.OptionFunc[HandlerOptions](func(cfg *HandlerOptions) {
		cfg.replacer = fn
	})
}
