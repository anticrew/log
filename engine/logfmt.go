package engine

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"
	"github.com/anticrew/log/internal/pool"
)

var _logFmtMarshalerPool = pool.NewPool(func() *LogFmtMarshaler {
	return &LogFmtMarshaler{}
})

type LogFmtMarshaler struct {
	out *buffer.Buffer
}

func NewLogFmtMarshaler(out *buffer.Buffer) *LogFmtMarshaler {
	j := _logFmtMarshalerPool.Get()
	j.out = out

	return j
}

func (j *LogFmtMarshaler) Dispose() {
	j.out = nil
	_logFmtMarshalerPool.Put(j)
}

func (j *LogFmtMarshaler) Marshal(r *log.Record) error {
	j.writeKey(TimeKey)
	j.out.WriteTime(r.Time, TimeLayout)

	j.writeKey(LevelKey)
	j.writeString(r.Level.String())

	j.writeKey(MessageKey)
	j.writeString(r.Message)

	var err error
	r.Attrs.Range(func(attr log.Attr) bool {
		err = errors.Join(j.writeAttr(attr))
		return false
	})

	if err != nil {
		j.out.Reset()
		return err
	}

	return nil
}

func (j *LogFmtMarshaler) writeKey(key string) {
	if isDirty(j.out, 0) {
		j.out.WriteString(" ")
	}

	j.out.WriteString(key).WriteByte('=')
}

func (j *LogFmtMarshaler) writeAttr(a log.Attr) error {
	j.writeKey(a.Key)

	switch a.Value.Kind() {
	case log.KindAny:
		j.out.WriteByte(buffer.Quote)

		err := json.NewEncoder(j.out).Encode(a.Value.Any())
		if err != nil {
			return err
		}

		j.out.CutSuffix([]byte{'\n'})
		j.out.WriteByte(buffer.Quote)

	case log.KindBool:
		j.out.WriteBool(a.Value.Bool())

	case log.KindDuration:
		j.out.WriteString(a.Value.Duration().String())

	case log.KindFloat64:
		j.out.WriteFloat64(a.Value.Float64(), 64)

	case log.KindInt64:
		j.out.WriteInt64(a.Value.Int64())

	case log.KindUint64:
		j.out.WriteUint64(a.Value.Uint64())

	case log.KindString:
		j.writeString(a.Value.String())

	case log.KindTime:
		j.out.WriteTime(a.Value.Time(), TimeLayout)
	}

	return nil
}

func (j *LogFmtMarshaler) writeString(s string) {
	if strings.ContainsRune(s, ' ') {
		j.out.WithQuotes().WriteString(s)
		return
	}

	j.out.WriteString(s)
}
