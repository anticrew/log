package log

import (
	"slices"
	"time"

	"github.com/anticrew/log/internal/pool"
)

const (
	ErrorKey  string = "error"
	CallerKey string = "caller"
)

type Attr struct {
	Key   string
	Value Value
}

func Err(err error) Attr {
	if err == nil {
		return String(ErrorKey, "nil")
	}

	return String(ErrorKey, err.Error())
}

func String(key, value string) Attr {
	return Attr{
		Key:   key,
		Value: StringValue(value),
	}
}

func Int64(key string, value int64) Attr {
	return Attr{
		Key:   key,
		Value: Int64Value(value),
	}
}

func Int(key string, value int) Attr {
	return Attr{
		Key:   key,
		Value: IntValue(value),
	}
}

func Uint64(key string, v uint64) Attr {
	return Attr{
		Key:   key,
		Value: Uint64Value(v),
	}
}

func Float64(key string, v float64) Attr {
	return Attr{
		Key:   key,
		Value: Float64Value(v),
	}
}

func Bool(key string, v bool) Attr {
	return Attr{
		Key:   key,
		Value: BoolValue(v),
	}
}

func Time(key string, v time.Time) Attr {
	return Attr{
		Key:   key,
		Value: TimeValue(v),
	}
}

func Duration(key string, v time.Duration) Attr {
	return Attr{
		Key:   key,
		Value: DurationValue(v),
	}
}

func Any(key string, a any) Attr {
	return Attr{
		Key:   key,
		Value: AnyValue(a),
	}
}

var _attrsPool = pool.NewPool(func() *Attrs {
	return &Attrs{}
})

type Attrs struct {
	attrs   []Attr
	storage [64]Attr
}

func NewAttrs() *Attrs {
	a := _attrsPool.Get()
	a.attrs = a.storage[:0]

	return a
}

func (a *Attrs) Range(f func(a Attr) bool) {
	for _, attr := range a.attrs {
		if f(attr) {
			break
		}
	}
}

func (a *Attrs) Search(key string) (Attr, bool) {
	ix, ok := slices.BinarySearchFunc(a.attrs, Attr{Key: key}, func(e Attr, t Attr) int {
		if e.Key < t.Key {
			return -1
		}
		if e.Key > t.Key {
			return 1
		}

		return 0
	})
	if ok {
		return a.attrs[ix], true
	}

	return Attr{}, false
}

func (a *Attrs) Append(attrs ...Attr) *Attrs {
	for _, attr := range attrs {
		a.appendSorted(attr)
	}

	return a
}

func (a *Attrs) Len() int {
	return len(a.attrs)
}

func (a *Attrs) Clone() *Attrs {
	return NewAttrs().Append(a.attrs...)
}

func (a *Attrs) Dispose() {
	a.attrs = nil
	_attrsPool.Put(a)
}

func (a *Attrs) appendSorted(attr Attr) {
	ix, ok := slices.BinarySearchFunc(a.attrs, attr, func(e, t Attr) int {
		if e.Key < t.Key {
			return -1
		}
		if e.Key > t.Key {
			return 1
		}

		return 0
	})
	if ok {
		a.attrs[ix] = attr
		return
	}

	a.attrs = slices.Insert(a.attrs, ix, attr)
}
