package engine

import (
	"encoding/json"
	"time"

	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"
	"github.com/anticrew/log/internal/pool"
)

var _prettyMarshalerPool = pool.NewPool(func() *PrettyMarshaler {
	return &PrettyMarshaler{}
})

type PrettyMarshaler struct {
	out *buffer.Buffer
}

func NewPrettyMarshaler(out *buffer.Buffer) *PrettyMarshaler {
	j := _prettyMarshalerPool.Get()
	j.out = out

	return j
}

func (j *PrettyMarshaler) Dispose() {
	j.out = nil
	_prettyMarshalerPool.Put(j)
}

func (j *PrettyMarshaler) Marshal(r *log.Record) error {
	j.writeTime(r.Time)
	j.out.WriteByte('\t')

	r.Attrs.Range(j.writeCaller)
	j.out.WriteByte('\n')

	j.writeLevel(r.Level)
	j.writeMessage(r.Message)
	j.out.WriteByte('\n')

	r.Attrs.Range(j.writeAttr)

	return nil
}

func (j *PrettyMarshaler) writeTime(t time.Time) {
	j.out.WriteTime(t, "[ 2006-01-02 15:04:05 ]")
}

func (j *PrettyMarshaler) writeCaller(a log.Attr) bool {
	if a.Key != log.CallerKey {
		return false
	}

	j.out.WriteString(a.Value.String())
	return true
}

func (j *PrettyMarshaler) writeLevel(l log.Level) {
	j.out.WriteByte(' ').WriteString(l.String()).WriteByte(' ')
}

func (j *PrettyMarshaler) writeMessage(msg string) {
	j.out.WriteByte(' ').WriteString(msg)
}

func (j *PrettyMarshaler) writeAttr(a log.Attr) bool {
	if a.Key == log.CallerKey {
		return false
	}

	j.out.WriteString("  ")
	j.out.WriteString(a.Key)

	j.out.WriteString(": ")
	err := j.writeAttrValue(a.Value)

	j.out.WriteByte('\n')

	if err != nil {
		j.out.WriteString("    ")
		j.out.WriteString(err.Error())
	}

	return false
}

func (j *PrettyMarshaler) writeAttrValue(v log.Value) error {
	switch v.Kind() {
	case log.KindAny:
		j.out.WriteByte(buffer.Quote)

		err := json.NewEncoder(j.out).Encode(v.Any())
		if err != nil {
			return err
		}

		j.out.CutSuffix([]byte{'\n'})
		j.out.WriteByte(buffer.Quote)

	case log.KindBool:
		j.out.WriteBool(v.Bool())

	case log.KindDuration:
		j.out.WriteString(v.Duration().String())

	case log.KindFloat64:
		j.out.WriteFloat64(v.Float64(), 64)

	case log.KindInt64:
		j.out.WriteInt64(v.Int64())

	case log.KindUint64:
		j.out.WriteUint64(v.Uint64())

	case log.KindString:
		j.out.WriteString(v.String())

	case log.KindTime:
		j.out.WriteTime(v.Time(), TimeLayout)
	}

	return nil
}
