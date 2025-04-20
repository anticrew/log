package engine

import (
	"io"

	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"
)

type Engine struct {
	out  io.Writer
	mode log.PrintMode
}

func NewEngine(out io.Writer, mode log.PrintMode) *Engine {
	return &Engine{
		out:  out,
		mode: mode,
	}
}

func (e *Engine) Write(r *log.Record) error {
	buf := buffer.New()
	defer buf.Dispose()

	var m Marshaler
	switch e.mode {
	case log.PrintModeLogFmt:
		m = NewLogFmtMarshaler(buf)
	case log.PrintModeJson:
		m = NewJsonMarshaler(buf)
	case log.PrintModePretty:
		m = NewPrettyMarshaler(buf)
	}

	defer m.Dispose()

	if err := m.Marshal(r); err != nil {
		return err
	}

	buf.WriteByte('\n')
	if _, err := buf.WriteTo(e.out); err != nil {
		return err
	}

	return nil
}
