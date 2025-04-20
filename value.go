package log

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"unsafe"
)

type Kind int

const (
	KindAny Kind = iota
	KindBool
	KindDuration
	KindFloat64
	KindInt64
	KindUint64
	KindString
	KindTime
)

type Value struct {
	num  uint64
	any  any
	kind Kind
}

func (v Value) Kind() Kind {
	return v.kind
}

func (v Value) Any() any {
	switch v.kind {
	case KindAny, KindTime, KindDuration:
		return v.any
	case KindFloat64, KindInt64, KindUint64:
		return v.num
	case KindBool:
		return v.num == 1
	case KindString:
		return v.unpackString()
	default:
		return fmt.Errorf(`%w "%d"`, ErrUnknownKind, v.kind)
	}
}

func (v Value) String() string {
	if v.kind == KindString {
		return v.unpackString()
	}

	return string(v.Bytes())
}

func (v Value) Bytes() []byte {
	var dst []byte

	switch v.Kind() {
	case KindString:
		return append(dst, v.unpackString()...)
	case KindInt64:
		return strconv.AppendInt(dst, int64(v.num), 10) //nolint: gosec // always int64, not uint64
	case KindUint64:
		return strconv.AppendUint(dst, v.num, 10)
	case KindFloat64:
		return strconv.AppendFloat(dst, v.unpackFloat64(), 'g', -1, 64)
	case KindBool:
		return strconv.AppendBool(dst, v.unpackBool())
	case KindDuration:
		d, _ := v.any.(time.Duration) //nolint: errcheck // always time.Duration
		return append(dst, d.String()...)
	case KindTime:
		t, _ := v.any.(time.Time) //nolint: errcheck // always time.Time
		return append(dst, t.String()...)
	case KindAny:
		return fmt.Append(dst, v.any)
	default:
		panic(fmt.Sprintf("bad kind: %d", v.kind))
	}
}

func (v Value) Bool() bool {
	return v.unpackBool()
}

func (v Value) Time() time.Time {
	t, _ := v.any.(time.Time) //nolint: errcheck // always time.Time
	return t
}

func (v Value) Duration() time.Duration {
	d, _ := v.any.(time.Duration) //nolint: errcheck // always time.Duration
	return d
}

func (v Value) Float64() float64 {
	return v.unpackFloat64()
}

func (v Value) Int64() int64 {
	return int64(v.num) //nolint: gosec // always int64, not uint64
}

func (v Value) Uint64() uint64 {
	return v.num
}

func (v Value) unpackString() string {
	return unsafe.String(v.any.(*byte), v.num) //nolint: errcheck // always *byte
}

func (v Value) unpackFloat64() float64 {
	return math.Float64frombits(v.num)
}

func (v Value) unpackBool() bool {
	return v.num == 1
}

func StringValue(v string) Value {
	return Value{
		num:  uint64(len(v)),
		any:  unsafe.StringData(v),
		kind: KindString,
	}
}

func Int64Value(v int64) Value {
	return Value{
		num:  uint64(v), //nolint: gosec // will be converted again
		kind: KindInt64,
	}
}

func IntValue(v int) Value {
	return Value{
		num:  uint64(v), //nolint: gosec // will be converted again
		kind: KindInt64,
	}
}

func Uint64Value(v uint64) Value {
	return Value{
		num:  v,
		kind: KindUint64,
	}
}

func Float64Value(v float64) Value {
	return Value{
		num:  math.Float64bits(v),
		kind: KindFloat64,
	}
}

func BoolValue(v bool) Value {
	var i uint64
	if v {
		i = 1
	}

	return Value{
		num:  i,
		kind: KindBool,
	}
}

func TimeValue(v time.Time) Value {
	return Value{
		any:  v,
		kind: KindTime,
	}
}

func DurationValue(v time.Duration) Value {
	return Value{
		any:  v,
		kind: KindDuration,
	}
}

func AnyValue(v any) Value {
	switch v := v.(type) {
	case Value:
		return v
	case string:
		return StringValue(v)
	case int:
		return Int64Value(int64(v))
	case uint:
		return Uint64Value(uint64(v))
	case int64:
		return Int64Value(v)
	case uint64:
		return Uint64Value(v)
	case bool:
		return BoolValue(v)
	case time.Duration:
		return DurationValue(v)
	case time.Time:
		return TimeValue(v)
	case uint8:
		return Uint64Value(uint64(v))
	case uint16:
		return Uint64Value(uint64(v))
	case uint32:
		return Uint64Value(uint64(v))
	case uintptr:
		return Uint64Value(uint64(v))
	case int8:
		return Int64Value(int64(v))
	case int16:
		return Int64Value(int64(v))
	case int32:
		return Int64Value(int64(v))
	case float64:
		return Float64Value(v)
	case float32:
		return Float64Value(float64(v))
	default:
		return Value{
			any:  v,
			kind: KindAny,
		}
	}
}
