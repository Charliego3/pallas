package logger

import (
	"go.uber.org/zap/buffer"
	"sync"
)

type state struct {
	buf     *buffer.Buffer
	freeBuf bool           // should buf be freed?
	sep     string         // separator to write before next key
	prefix  *buffer.Buffer // for text: key prefix
	groups  *[]string      // pool-allocated slice of active groups, for ReplaceAttr
}

var groupPool = sync.Pool{New: func() any {
	s := make([]string, 0, 10)
	return &s
}}

func (h *TextHandler) newState(buf *buffer.Buffer, freeBuf bool, sep string, prefix *buffer.Buffer) *state {
	return &state{buf: buf, freeBuf: freeBuf, sep: sep, prefix: prefix}
}
