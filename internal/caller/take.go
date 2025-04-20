package caller

import (
	"github.com/anticrew/log/internal/buffer"
)

func Take(skipCount int) (string, error) {
	c, dispose, err := Capture(skipCount + 1) // skip current Take call
	if err != nil {
		return "", err
	}

	defer dispose()

	buf := buffer.New()
	defer buf.Dispose()

	buf.WriteString(c.File).
		WriteByte(':').
		WriteInt64(int64(c.Line)).
		WriteByte(' ').
		WriteString(c.Function)

	return buf.String(), nil
}
