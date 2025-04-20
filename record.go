package log

import (
	"time"

	"github.com/anticrew/log/internal/pool"
)

type Record struct {
	Time    time.Time
	Level   Level
	Message string
	Attrs   *Attrs
}

func NewRecord(level Level, message string, attrs *Attrs) *Record {
	r := _recordPool.Get()
	r.Time = time.Now()
	r.Level = level
	r.Message = message
	r.Attrs = attrs
	return r
}

func (r *Record) Dispose() {
	r.Attrs.Dispose()
	r.Attrs = nil
	_recordPool.Put(r)
}

var _recordPool = pool.NewPool(func() *Record {
	return &Record{}
})
