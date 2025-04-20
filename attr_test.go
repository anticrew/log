package log

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	t.Parallel()

	type testCase struct {
		err           error
		expectedValue Value
	}

	tests := map[string]testCase{
		"nil": {
			err:           nil,
			expectedValue: AnyValue(nil),
		},
		"err": {
			err: assert.AnError,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			attr := Err(test.err)

			assert.Equal(t, ErrorKey, attr.Key)
			assert.Equal(t, KindString, attr.Value.kind)

			if test.err == nil {
				assert.Equal(t, "nil", attr.Value.unpackString())
			} else {
				assert.Equal(t, test.err.Error(), attr.Value.unpackString())
			}
		})
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value string
	}

	tests := map[string]testCase{
		"пустая строка": {
			key:   "key",
			value: "",
		},
		"непустая строка": {
			key:   "test",
			value: "value",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := String(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.Equal(t, test.value, attr.Value.String())
			assert.Equal(t, KindString, attr.Value.kind)
		})
	}
}

func TestInt64(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value int64
	}

	tests := map[string]testCase{
		"ноль": {
			key:   "key",
			value: 0,
		},
		"положительное число": {
			key:   "pos",
			value: 42,
		},
		"отрицательное число": {
			key:   "neg",
			value: -42,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Int64(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.Equal(t, test.value, attr.Value.Int64())
			assert.Equal(t, KindInt64, attr.Value.kind)
		})
	}
}

func TestTime(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value time.Time
	}

	now := time.Now()
	tests := map[string]testCase{
		"нулевое время": {
			key:   "empty",
			value: time.Time{},
		},
		"текущее время": {
			key:   "now",
			value: now,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Time(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.Equal(t, test.value, attr.Value.Time())
			assert.Equal(t, KindTime, attr.Value.kind)
		})
	}
}

func TestInt(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value int
	}

	tests := map[string]testCase{
		"ноль": {
			key:   "zero",
			value: 0,
		},
		"положительное": {
			key:   "pos",
			value: 100,
		},
		"отрицательное": {
			key:   "neg",
			value: -100,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Int(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.EqualValues(t, test.value, attr.Value.Int64())
			assert.Equal(t, KindInt64, attr.Value.kind)
		})
	}
}

func TestUint64(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value uint64
	}

	tests := map[string]testCase{
		"ноль": {
			key:   "zero",
			value: 0,
		},
		"максимальное": {
			key:   "max",
			value: ^uint64(0),
		},
		"обычное": {
			key:   "normal",
			value: 42,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Uint64(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.Equal(t, test.value, attr.Value.Uint64())
			assert.Equal(t, KindUint64, attr.Value.kind)
		})
	}
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value float64
	}

	tests := map[string]testCase{
		"ноль": {
			key:   "zero",
			value: 0.0,
		},
		"положительное": {
			key:   "pos",
			value: 3.14,
		},
		"отрицательное": {
			key:   "neg",
			value: -2.718,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Float64(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.InDelta(t, test.value, attr.Value.Float64(), 0.000001)
			assert.Equal(t, KindFloat64, attr.Value.kind)
		})
	}
}

func TestBool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value bool
	}

	tests := map[string]testCase{
		"true": {
			key:   "true",
			value: true,
		},
		"false": {
			key:   "false",
			value: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Bool(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.Equal(t, test.value, attr.Value.Bool())
			assert.Equal(t, KindBool, attr.Value.kind)
		})
	}
}

func TestDuration(t *testing.T) {
	t.Parallel()

	type testCase struct {
		key   string
		value time.Duration
	}

	tests := map[string]testCase{
		"ноль": {
			key:   "zero",
			value: 0,
		},
		"секунда": {
			key:   "sec",
			value: time.Second,
		},
		"отрицательная": {
			key:   "neg",
			value: -time.Hour,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			attr := Duration(test.key, test.value)

			assert.Equal(t, test.key, attr.Key)
			assert.Equal(t, test.value, attr.Value.Duration())
			assert.Equal(t, KindDuration, attr.Value.kind)
		})
	}
}

func TestAttrs_Append(t *testing.T) {
	t.Parallel()

	type testCase struct {
		attrs       *Attrs
		expectedLen int
	}

	const (
		duplicateKey      = "key"
		duplicateNewValue = "new-value"
	)

	tests := map[string]testCase{
		"empty": {
			attrs:       NewAttrs(),
			expectedLen: 1,
		},
		"append": {
			attrs:       NewAttrs().Append(Duration("duration", time.Second)),
			expectedLen: 2,
		},
		"rewrite": {
			attrs:       NewAttrs().Append(Int(duplicateKey, 42)),
			expectedLen: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer test.attrs.Dispose()

			test.attrs.Append(String(duplicateKey, duplicateNewValue))
			assert.Equal(t, test.expectedLen, test.attrs.Len())

			test.attrs.Range(func(a Attr) bool {
				if a.Key != duplicateKey {
					return false
				}

				assert.Equal(t, duplicateNewValue, a.Value.String())
				return true
			})
		})
	}
}

func TestAttrs_Range(t *testing.T) {
	t.Parallel()

	type testCase struct {
		attrs         *Attrs
		fn            func(a Attr) bool
		expectedCount int
	}

	var (
		intAttr    = Int("int", 42)
		stringAttr = String("string", "value")
	)

	testCases := map[string]testCase{
		"empty": {
			attrs:         NewAttrs(),
			fn:            func(a Attr) bool { return true },
			expectedCount: 0,
		},
		"one": {
			attrs:         NewAttrs().Append(stringAttr),
			fn:            func(a Attr) bool { return true },
			expectedCount: 1,
		},
		"two": {
			attrs:         NewAttrs().Append(intAttr, stringAttr),
			fn:            func(a Attr) bool { return a.Key == intAttr.Key },
			expectedCount: 1,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer test.attrs.Dispose()

			var actualCount int
			test.attrs.Range(func(a Attr) bool {
				actualCount++
				return test.fn(a)
			})

			assert.Equal(t, test.expectedCount, actualCount)
		})
	}
}

func TestAttrs_Search(t *testing.T) {
	t.Parallel()

	type testCase struct {
		attrs        *Attrs
		key          string
		expected     bool
		expectedAttr Attr
	}

	var (
		intAttr    = Int("int", 42)
		stringAttr = String("string", "value")
	)

	testCases := map[string]testCase{
		"empty": {
			attrs:    NewAttrs(),
			key:      "zero",
			expected: false,
		},
		"one": {
			attrs:        NewAttrs().Append(intAttr),
			key:          intAttr.Key,
			expected:     true,
			expectedAttr: intAttr,
		},
		"two": {
			attrs:        NewAttrs().Append(stringAttr, intAttr),
			key:          stringAttr.Key,
			expected:     true,
			expectedAttr: stringAttr,
		},
		"two-not-found": {
			attrs:    NewAttrs().Append(stringAttr, intAttr),
			key:      "key",
			expected: false,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer test.attrs.Dispose()

			attr, ok := test.attrs.Search(test.key)
			assert.Equal(t, test.expected, ok)
			assert.Equal(t, test.expectedAttr, attr)
		})
	}
}

func TestAttrs_Len(t *testing.T) {
	t.Parallel()

	type testCase struct {
		attrs       *Attrs
		expectedLen int
	}

	testCases := map[string]testCase{
		"empty": {
			attrs:       NewAttrs(),
			expectedLen: 0,
		},
		"one": {
			attrs:       NewAttrs().Append(String("string", "value")),
			expectedLen: 1,
		},
		"two": {
			attrs:       NewAttrs().Append(Int("int", 42), Duration("duration", time.Second*2)),
			expectedLen: 2,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer test.attrs.Dispose()

			assert.Equal(t, test.expectedLen, test.attrs.Len())
		})
	}
}

func TestAttrs_Clone(t *testing.T) {
	t.Parallel()

	attrs := NewAttrs()
	defer attrs.Dispose()

	attrs.Append(String("string", "value"))
	assert.Equal(t, 1, attrs.Len())

	attrsCopy := attrs.Clone()
	assert.Equal(t, attrs.Len(), attrsCopy.Len())
	assert.NotEqual(t, fmt.Sprintf("%p", attrs), fmt.Sprintf("%p", attrsCopy)) // сравниваем указатели

	attrsCopy.Append(Int("int", 42))
	assert.Equal(t, 2, attrsCopy.Len())
	assert.Equal(t, 1, attrs.Len())
}
