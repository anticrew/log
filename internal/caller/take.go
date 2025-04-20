package caller

import (
	"strings"

	"github.com/anticrew/go-x/xio"
)

// Take
// Возвращает строковое представление источника вызова в компактном формате: информация о файле сокращается до названия пакета и функции
func Take(skipCount int) (string, error) {
	c, dispose, err := Capture(skipCount + 1) // skip current Take call
	if err != nil {
		return "", err
	}

	defer dispose()

	buf := xio.NewBuffer()
	defer buf.Dispose()

	// Find the last separator.
	idx := strings.LastIndexByte(c.File, '/')
	if idx == -1 {
		appendFullPath(buf, c)
		return buf.String(), nil
	}

	// Find the penultimate separator.
	idx = strings.LastIndexByte(c.File[:idx], '/')
	if idx == -1 {
		appendFullPath(buf, c)
		return buf.String(), nil
	}

	buf.WriteString(c.File[idx+1:]).
		WriteByte(':').
		WriteInt64(int64(c.Line)).
		WriteByte(' ').
		WriteString(c.Function)

	return buf.String(), nil
}

func appendFullPath(buf *xio.Buffer, c *Caller) {
	buf.WriteString(c.File).
		WriteByte(':').
		WriteInt64(int64(c.Line)).
		WriteByte(' ').
		WriteString(c.Function)
}
