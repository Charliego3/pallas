package logger

import "time"

type textBuilder struct {
	*baseBuilder
}

func (b *textBuilder) appendTime(t time.Time) {
	if t.IsZero() {
		return
	}

	b.buf.WriteString(timeStyle.Render(t.Format(b.opts.timeFormat)))
	b.buf.WriteByte(b.componmentSep())
}
