package interpolation

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anticrew/log"

	"github.com/stretchr/testify/assert"
)

func Test_replace(t *testing.T) {
	t.Parallel()

	type testCase struct {
		src      string
		attrs    *log.Attrs
		expected string
		err      error
	}

	testCases := map[string]testCase{
		"no-interpolation-lat": {
			src:      "Hello world!",
			expected: "Hello world!",
		},
		"no-interpolation-cyr": {
			src:      "Привет мир!",
			expected: "Привет мир!",
		},
		"no-interpolation-with-bracket-lat": {
			src:      "Hello {} world!",
			expected: "Hello {} world!",
		},
		"no-interpolation-with-bracket-cyr": {
			src:      "Привет {} мир!",
			expected: "Привет {} мир!",
		},
		"no-attrs-lat": {
			src:      "Hello {{ name }}!",
			expected: "Hello {{ name }}!",
		},
		"no-attrs-cyr": {
			src:      "Привет {{ name }}!",
			expected: "Привет {{ name }}!",
		},
		"no-match-attrs-lat": {
			src:      "Hello {{ name }}!",
			expected: "Hello {{ name }}!",
			attrs:    log.NewAttrs().Append(log.String("world", "world")),
		},
		"no-match-attrs-cyr": {
			src:      "Привет {{ name }}!",
			expected: "Привет {{ name }}!",
			attrs:    log.NewAttrs().Append(log.String("world", "world")),
		},
		"string-attrs-lat": {
			src:      "Hello {{ name }}!",
			attrs:    log.NewAttrs().Append(log.String("name", "world")),
			expected: "Hello world!",
		},
		"string-attrs-cyr": {
			src:      "Привет {{ name }}!",
			attrs:    log.NewAttrs().Append(log.String("name", "мир")),
			expected: "Привет мир!",
		},
		"mixed-attrs-lat": {
			src: `Hello {{ name }}! There is more than {{ count }} test cases. ` +
				`Anywhere error can be occurred, for example "{{ err }}"`,
			attrs: log.NewAttrs().Append(
				log.String("name", "world"),
				log.Int("count", 7),
				log.Any("err", errors.New("fatal error")),
			),
			expected: `Hello world! There is more than 7 test cases. ` +
				`Anywhere error can be occurred, for example "fatal error"`,
		},
		"mixed-attrs-cyr": {
			src: `Привет {{ name }}! Здесь больше {{ count }} тест кейсов. ` +
				`Везде может возникнуть ошибка, к примеру "{{ err }}"`,
			attrs: log.NewAttrs().Append(
				log.String("name", "мир"),
				log.Int("count", 7),
				log.Any("err", errors.New("фатальная ошибка")),
			),
			expected: `Привет мир! Здесь больше 7 тест кейсов. Везде может возникнуть ошибка, ` +
				`к примеру "фатальная ошибка"`,
		},
		"joined-attrs-lat": {
			src: "Hello {{ name }}{{ char }}",
			attrs: log.NewAttrs().Append(
				log.String("name", "world"),
				log.String("char", "!"),
			),
			expected: "Hello world!",
		},
		"joined-attrs-cyr": {
			src: "Привет {{ name }}{{ char }}",
			attrs: log.NewAttrs().Append(
				log.String("name", "мир"),
				log.String("char", "!"),
			),
			expected: "Привет мир!",
		},
		"too-much-attrs": {
			src: "Hello {{ name }}!",
			attrs: log.NewAttrs().Append(
				log.String("name", "world"),
				log.Any("random", 'z'),
			),
			expected: "Hello world!",
		},
		"err-indirect-open": {
			src:   "Hello {{ name {{ key!",
			attrs: log.NewAttrs().Append(log.String("name", "world")),
			err:   ErrIndirectOpenKey,
		},
		"err-indirect-close": {
			src:   "Hello name }}!",
			attrs: log.NewAttrs().Append(log.String("name", "world")),
			err:   ErrIndirectCloseKey,
		},
		"err-empty-key": {
			src:   "Hello {{}}!",
			attrs: log.NewAttrs().Append(log.String("name", "world")),
			err:   ErrEmptyKey,
		},
		"err-unclosed-key": {
			src:   "Hello {{ key",
			attrs: log.NewAttrs().Append(log.String("name", "world")),
			err:   ErrUnclosedKey,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := newReplacer(test.src, test.attrs)
			defer r.dispose()

			actual, err := r.replace()
			if test.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, test.err)
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_next(t *testing.T) {
	t.Parallel()

	type testCase struct {
		data  string
		i     int
		value byte
		ok    bool
	}

	testCases := map[string]testCase{
		"empty": {},
		"start": {
			data:  "abc",
			i:     -1,
			value: 'a',
			ok:    true,
		},
		"middle": {
			data:  "abc",
			i:     0,
			value: 'b',
			ok:    true,
		},
		"end": {
			data:  "abc",
			i:     1,
			value: 'c',
			ok:    true,
		},
		"out-of-range": {
			data: "abc",
			i:    2,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := newReplacer(test.data, nil)
			defer r.dispose()

			r.i = test.i

			v, ok := r.next()
			assert.Equal(t, test.value, v)
			assert.Equal(t, test.ok, ok)
		})
	}
}

var (
	benchmarkSource = `Привет {name}! Здесь больше {count} тест кейсов. Везде может возникнуть ошибка, к примеру "{err}"`
	benchmarkAttrs  = log.NewAttrs().Append(
		log.String("name", "мир"),
		log.Int("count", 7),
		log.Any("err", errors.New("фатальная ошибка")),
	)
)

func Benchmark_replace(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		r := newReplacer(benchmarkSource, benchmarkAttrs)

		s, err := r.replace()
		_ = s
		_ = err

		r.dispose()
	}
}

func Benchmark_Replacer(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		replaces := make([]string, 0, 64)
		benchmarkAttrs.Range(func(a log.Attr) bool {
			replaces = append(replaces, strings.Join([]string{"{", a.Key, "}"}, ""), a.Value.String())

			return false
		})

		r := strings.NewReplacer(replaces...)
		s := r.Replace(benchmarkSource)
		_ = s
	}
}
