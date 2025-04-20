package engine

import (
	"encoding/json"
	"errors"

	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"
	"github.com/anticrew/log/internal/maps"
	"github.com/anticrew/log/internal/pool"
)

var _jsonMarshalerPool = pool.NewPool(func() *JsonMarshaler {
	return &JsonMarshaler{
		knownKeys: make(map[string]struct{}, 64),
	}
})

type JsonMarshaler struct {
	knownKeys map[string]struct{}
	out       *buffer.Buffer
}

func NewJsonMarshaler(out *buffer.Buffer) *JsonMarshaler {
	j := _jsonMarshalerPool.Get()
	j.out = out

	return j
}

func (j *JsonMarshaler) Dispose() {
	j.knownKeys = maps.Clear(j.knownKeys)
	_jsonMarshalerPool.Put(j)
}

func (j *JsonMarshaler) Marshal(r *log.Record) error {
	const (
		openBracket  = '{'
		closeBracket = '}'
	)

	j.out.WriteByte(openBracket)

	// никогда не завершается ошибкой, т. к. ключей в knownKeys еще нет
	_ = j.writeKey(TimeKey)

	j.out.WithQuotes().WriteTime(r.Time, TimeLayout)

	if err := j.writeKey(LevelKey); err != nil {
		return err
	}

	j.out.WithQuotes().WriteString(r.Level.String())

	if err := j.writeKey(MessageKey); err != nil {
		return err
	}

	j.out.WithQuotes().WriteString(r.Message)

	var err error
	r.Attrs.Range(func(attr log.Attr) bool {
		err = errors.Join(j.writeAttr(attr))
		return false
	})

	if err != nil {
		j.out.Reset()
		return err
	}

	j.out.WriteByte(closeBracket)
	return nil
}

func (j *JsonMarshaler) writeKey(key string) (err error) {
	if _, ok := j.knownKeys[key]; ok {
		return ErrKeyExists
	}
	j.knownKeys[key] = struct{}{}

	// 1 разрешает наличие открывающей скобки { в буфере
	if isDirty(j.out, 1) {
		j.out.WriteString(",")
	}

	j.out.WithQuotes().WriteString(key).WriteByte(':')

	return nil
}

func (j *JsonMarshaler) writeAttr(a log.Attr) error {
	if err := j.writeKey(a.Key); err != nil {
		return err
	}

	switch a.Value.Kind() {
	case log.KindAny:
		err := json.NewEncoder(j.out).Encode(a.Value.Any())
		if err != nil {
			return err
		}

		j.out.CutSuffix([]byte{'\n'})

	case log.KindBool:
		j.out.WriteBool(a.Value.Bool())

	case log.KindDuration:
		j.out.WithQuotes().WriteString(a.Value.Duration().String())

	case log.KindFloat64:
		j.out.WriteFloat64(a.Value.Float64(), 64)

	case log.KindInt64:
		j.out.WriteInt64(a.Value.Int64())

	case log.KindUint64:
		j.out.WriteUint64(a.Value.Uint64())

	case log.KindString:
		j.out.WithQuotes().WriteString(a.Value.String())

	case log.KindTime:
		j.out.WithQuotes().WriteTime(a.Value.Time(), TimeLayout)
	}

	return nil
}
