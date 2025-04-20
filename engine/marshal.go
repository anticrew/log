package engine

import (
	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"
)

type Marshaler interface {
	Marshal(r *log.Record) error
	Dispose()
}

func isDirty(b *buffer.Buffer, maxLen int) bool {
	return b.Len() > maxLen
}
