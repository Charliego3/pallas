package logger

import "time"

type jsonBuilder struct {
	*baseBuilder
}

func (b *jsonBuilder) appendTime(t time.Time) {
	if t.IsZero() {
		return
	}

	b.buf.WriteString(timeStyle.Render(t.Format(b.opts.timeFormat)))
	b.buf.WriteByte(b.componmentSep())
}

func (b *jsonBuilder) start() {
	b.buf.WriteByte('{')
}

func (b *jsonBuilder) end() {
	b.buf.WriteByte('}')
}
