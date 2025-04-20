package buffer

import (
	"io"
	"strconv"
	"time"

	"github.com/anticrew/log/internal/pool"
)

const (
	Quote = '"'

	defaultSize      = 1024
	effectiveGrowMin = 32
)

var _pool = pool.NewPool(func() *Buffer {
	return &Buffer{
		buf: make([]byte, 0, defaultSize),
	}
})

type Buffer struct {
	buf    []byte
	quotes bool
}

func New() *Buffer {
	return _pool.Get()
}

func (b *Buffer) WithQuotes() *Buffer {
	b.quotes = true
	return b
}

func (b *Buffer) WriteByte(v byte) *Buffer {
	b.writeQuote()
	b.buf = append(b.buf, v)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteBytes(v []byte) *Buffer {
	b.writeQuote()
	b.buf = append(b.buf, v...)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteString(s string) *Buffer {
	b.writeQuote()
	b.buf = append(b.buf, s...)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteInt64(i int64) *Buffer {
	b.writeQuote()
	b.buf = strconv.AppendInt(b.buf, i, 10)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteUint64(i uint64) *Buffer {
	b.writeQuote()
	b.buf = strconv.AppendUint(b.buf, i, 10)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteFloat64(f float64, bitSize int) *Buffer {
	b.writeQuote()
	b.buf = strconv.AppendFloat(b.buf, f, 'f', -1, bitSize)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteBool(v bool) *Buffer {
	b.writeQuote()
	b.buf = strconv.AppendBool(b.buf, v)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) WriteTime(t time.Time, layout string) *Buffer {
	b.writeQuote()
	b.buf = t.AppendFormat(b.buf, layout)
	b.writeQuote()
	b.quotes = false

	return b
}

func (b *Buffer) Write(v []byte) (int, error) {
	b.WriteBytes(v)
	return len(v), nil
}

func (b *Buffer) Len() int {
	return len(b.buf)
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

func (b *Buffer) String() string {
	return string(b.buf)
}

func (b *Buffer) Reset() {
	b.quotes = false
	b.buf = b.buf[:0]
}

func (b *Buffer) Dispose() {
	disposeBuffer(b)
}

func (b *Buffer) CutSuffix(suffix []byte) *Buffer {
	var (
		bufLen    = len(b.buf)
		suffixLen = len(suffix)
	)
	if suffixLen > bufLen {
		return b
	}

	for i := range suffixLen {
		bufIndex := bufLen - i - 1
		sufIndex := suffixLen - i - 1

		if b.buf[bufIndex] != suffix[sufIndex] {
			return b
		}
	}

	b.buf = b.buf[:bufLen-suffixLen]

	return b
}

func (b *Buffer) writeQuote() {
	if !b.quotes {
		return
	}

	b.buf = append(b.buf, Quote)
}

func (b *Buffer) WriteTo(out io.Writer) (int64, error) {
	n, err := out.Write(b.buf)
	return int64(n), err
}

func disposeBuffer(b *Buffer) {
	b.Reset()
	_pool.Put(b)
}
