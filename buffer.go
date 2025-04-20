package log

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/anticrew/go-x/xio"
)

type BufferWriter struct {
	mx     *sync.Mutex
	buffer *xio.Buffer
	out    io.Writer
}

func NewBufferWriter(ctx context.Context, delay time.Duration, out io.Writer) *BufferWriter {
	if delay <= 0 {
		delay = 100 * time.Millisecond
	}

	b := &BufferWriter{
		mx:     &sync.Mutex{},
		buffer: xio.NewBuffer(),
		out:    out,
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				_ = b.Close()
				return
			case <-time.After(delay):
				b.flush()
			}
		}
	}(ctx)

	return b
}

func (b *BufferWriter) Write(p []byte) (n int, err error) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.buffer == nil {
		return 0, io.ErrClosedPipe
	}

	return b.buffer.Write(p)
}

func (b *BufferWriter) Close() error {
	b.flush()
	b.buffer.Dispose()
	b.buffer = nil

	return nil
}

func (b *BufferWriter) flush() {
	b.mx.Lock()
	defer b.mx.Unlock()

	if b.buffer == nil || b.out == nil {
		return
	}

	defer b.buffer.Reset()
	if _, err := b.out.Write(b.buffer.Bytes()); err != nil {
		return
	}
}
