package interpolation

import (
	"fmt"

	"github.com/anticrew/log/internal/pool"

	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"
)

type parseOp int

const (
	opUnknown parseOp = iota
	opChar
	opKeyOpen
	opKeyName
	opKeyClose
)

var _replacerPool = pool.NewPool(func() *replacer {
	return &replacer{}
})

type replacer struct {
	op      parseOp
	opStart int
	i       int

	keyOpenWhitespaces  int
	keyCloseWhitespaces int

	src   string
	attrs *log.Attrs

	buf *buffer.Buffer
}

func newReplacer(src string, attrs *log.Attrs) *replacer {
	r := _replacerPool.Get()
	r.op = opChar
	r.opStart = -1
	r.i = 0

	r.keyOpenWhitespaces = 0
	r.keyCloseWhitespaces = 0

	r.src = src
	r.attrs = attrs

	r.buf = buffer.New()

	return r
}

func (r *replacer) dispose() {
	r.buf.Dispose()
	r.buf = nil

	r.src = ""
	r.attrs = nil

	_replacerPool.Put(r)
}

// func (r *replacer) replace(src string, attrs *xlog.Attrs) (string, error) {
func (r *replacer) replace() (string, error) {
	if r.attrs == nil || r.attrs.Len() == 0 {
		return r.src, nil
	}

	if err := r.run(); err != nil {
		return "", err
	}

	if r.op != opChar && r.op != opKeyClose {
		return "", fmt.Errorf("%w at %d", ErrUnclosedKey, r.opStart)
	}

	if r.opStart > -1 && len(r.src)-r.opStart > 0 {
		r.buf.WriteString(r.src[r.opStart:])
	}

	return r.buf.String(), nil
}

func (r *replacer) run() error {
	for r.i = 0; r.i < len(r.src); r.i++ {
		if r.processWhitespaces() {
			continue
		}

		if ok, err := r.processOpenKey(); ok {
			if err != nil {
				return err
			}

			continue
		}

		if ok, err := r.processCloseKey(); ok {
			if err != nil {
				return err
			}

			continue
		}

		if r.op == opKeyOpen {
			r.op = opKeyName
			continue
		}

		if r.op == opKeyName {
			continue
		}

		r.op = opChar
		r.opStart = -1

		r.buf.WriteByte(r.src[r.i])
	}

	return nil
}

func (r *replacer) processWhitespaces() bool {
	v := r.src[r.i]
	if v != ' ' && v != '\t' && v != '\n' && v != '\r' {
		return false
	}

	//nolint: exhaustive // non influencing cases: interpolation.opUnknown, interpolation.opChar, interpolation.opKeyClose
	switch r.op {
	case opKeyOpen:
		r.keyOpenWhitespaces++

	case opKeyName:
		r.keyCloseWhitespaces++

	default:
		return false
	}

	return true
}

func (r *replacer) processOpenKey() (bool, error) {
	var (
		current  = r.src[r.i]
		next, ok = r.next()
	)

	if current != '{' || !ok || next != '{' {
		return false, nil
	}

	if r.op != opChar && r.op != opKeyClose {
		return true, fmt.Errorf("%w at %d", ErrIndirectOpenKey, r.i)
	}

	r.keyOpenWhitespaces = 0
	r.keyCloseWhitespaces = 0

	r.op = opKeyOpen
	r.opStart = r.i
	r.i++ // skip next {

	return true, nil
}

func (r *replacer) processCloseKey() (bool, error) {
	var (
		current  = r.src[r.i]
		next, ok = r.next()
	)

	if current != '}' || !ok || next != '}' {
		return false, nil
	}

	if r.op != opKeyName {
		if r.op == opKeyOpen {
			return true, fmt.Errorf("%w at %d", ErrEmptyKey, r.i)
		}

		return true, fmt.Errorf("%w at %d", ErrIndirectCloseKey, r.i)
	}

	// opStart - index of first byte in {{
	// keyOpenWhitespaces - count of whitespaces between {{ and first byte in key name
	// +2 skips {{
	key := r.src[r.opStart+r.keyOpenWhitespaces+2 : r.i-r.keyCloseWhitespaces]

	attr, ok := r.attrs.Search(key) //nolint: govet // shadow ok is normal
	if !ok {
		// missed attr: write origin placeholder without changes
		r.buf.WriteString(r.src[r.opStart : r.i+1])
		r.op = opChar
		return true, nil
	}

	r.op = opKeyClose
	r.opStart = -1
	r.i++ // skip next }

	r.buf.WriteString(attr.Value.String())
	return true, nil
}

func (r *replacer) next() (byte, bool) {
	if j := r.i + 1; j < len(r.src) {
		return r.src[j], true
	}

	return 0, false
}
