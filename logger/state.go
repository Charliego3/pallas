package logger

import (
	"github.com/charliego3/mspp/container"
	"sync"
)

type state struct {
	buf     *container.Buffer
	freeBuf bool              // should buf be freed?
	sep     string            // separator to write before next key
	prefix  *container.Buffer // for text: key prefix
	groups  *[]string         // pool-allocated slice of active groups, for ReplaceAttr
}

var groupPool = sync.Pool{New: func() any {
	s := make([]string, 0, 10)
	return &s
}}

func (h *TextHandler) newState(buf *container.Buffer, freeBuf bool, sep string, prefix *container.Buffer) *state {
	s := &state{buf: buf, freeBuf: freeBuf, sep: sep, prefix: prefix}
	if h.options.replacer != nil {
		s.groups = groupPool.Get().(*[]string)
		*s.groups = append(*s.groups, h.groups[:h.nOpenGroups]...)
	}
	return s
}

func (s *state) free() {
	if s.freeBuf {
		s.buf.Free()
	}
	if gs := s.groups; gs != nil {
		*gs = (*gs)[:0]
		groupPool.Put(gs)
	}
}
