package logger

import (
	"github.com/charliego3/mspp/container"
	"sync"
	"time"
)

const (
	textComponmentSep = '='
	jsonComponmentSep = ':'

	textAttrSep = ' '
	jsonAttrSep = ','
)

var groupPool = sync.Pool{New: func() any {
	s := make([]string, 0, 10)
	return &s
}}

type Builder interface {
	start()
	end()
	appendTime(time.Time)
	componmentSep() byte
	attrSep() byte
	free()
}

type baseBuilder struct {
	opts    *HandlerOptions
	buf     *container.Buffer
	freeBuf bool              // should buf be freed?
	sep     string            // separator to write before next key
	prefix  *container.Buffer // for text: key prefix
	groups  *[]string         // pool-allocated slice of active groups, for ReplaceAttr
	json    bool
}

func (b *baseBuilder) componmentSep() byte {
	if b.json {
		return jsonComponmentSep
	}
	return textComponmentSep
}

func (b *baseBuilder) attrSep() byte {
	if b.json {
		return jsonAttrSep
	}
	return textAttrSep
}

func (h *baseHandler) newBuilder(buf *container.Buffer, freeBuf bool, sep string, prefix *container.Buffer) Builder {
	builder := &baseBuilder{buf: buf, freeBuf: freeBuf, sep: sep, prefix: prefix, json: h.json}
	if h.options.replacer != nil {
		builder.groups = groupPool.Get().(*[]string)
		*builder.groups = append(*builder.groups, h.groups[:h.nOpenGroups]...)
	}
	if h.json {
		return &jsonBuilder{baseBuilder: builder}
	}
	return nil
}

func (b *baseBuilder) free() {
	if b.freeBuf {
		b.buf.Free()
	}
	if gs := b.groups; gs != nil {
		*gs = (*gs)[:0]
		groupPool.Put(gs)
	}
}
