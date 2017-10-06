package human

import (
	"bytes"
	"io"
)

type FlushableBuffer struct {
	*bytes.Buffer
	stream io.Writer
}

func (b *FlushableBuffer) Flush() (n int, err error) {
	return b.stream.Write(b.Bytes())
}

func NewFlushableBuffer(stream io.Writer) *FlushableBuffer {
	return &FlushableBuffer{
		Buffer: bytes.NewBufferString(""),
		stream: stream,
	}
}
