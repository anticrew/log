package log

import (
	"fmt"
	"strconv"

	"github.com/anticrew/log/internal/maps"
)

const (
	PrintModePretty PrintMode = iota
	PrintModeJson
	PrintModeLogFmt
)

var (
	_printModeNames = map[PrintMode]string{
		PrintModePretty: "PRETTY",
		PrintModeJson:   "JSON",
		PrintModeLogFmt: "TEXT",
	}

	_printModeByName = maps.Invert(_printModeNames)
)

func SetPrintModeName(printMode PrintMode, name string) {
	_printModeNames[printMode] = name
	_printModeByName = maps.Invert(_printModeNames)
}

type PrintMode int

func (p PrintMode) String() string {
	s, ok := _printModeNames[p]
	if ok {
		return s
	}

	return fmt.Sprintf("PrintMode<%d>", p)
}

func (p PrintMode) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *PrintMode) UnmarshalText(text []byte) error {
	return p.parse(string(text))
}

func (p PrintMode) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, p.String()), nil
}

func (p *PrintMode) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return p.parse(s)
}

func (p *PrintMode) parse(name string) error {
	var ok bool
	*p, ok = _printModeByName[name]
	if !ok {
		return fmt.Errorf(`%w "%s"`, ErrUnknownPrintMode, name)
	}

	return nil
}
