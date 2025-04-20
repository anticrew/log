package log

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/anticrew/log/internal/maps"
)

const (
	LevelTrace Level = -8
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
	LevelFatal Level = 12
)

var (
	_levelNames = map[Level]string{
		LevelTrace: "TRACE",
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
		LevelFatal: "FATAL",
	}

	_levelByName = maps.Invert(_levelNames)
)

func SetLevelName(level Level, name string) {
	_levelNames[level] = name
	_levelByName = maps.Invert(_levelNames)
}

type Level slog.Level

func (l Level) Level() Level {
	return l
}

func (l Level) String() string {
	s, ok := _levelNames[l]
	if ok {
		return s
	}

	return fmt.Sprintf("Level<%d>", l)
}

func (l Level) MarshalText() (text []byte, err error) {
	return []byte(l.String()), nil
}

func (l *Level) UnmarshalText(data []byte) error {
	return l.parse(string(data))
}

func (l Level) MarshalJSON() ([]byte, error) {
	return strconv.AppendQuote(nil, l.String()), nil
}

func (l *Level) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	return l.parse(s)
}

func (l *Level) parse(name string) error {
	var ok bool
	*l, ok = _levelByName[name]
	if !ok {
		return fmt.Errorf(`%w "%s"`, ErrUnknownLevel, name)
	}

	return nil
}

type Leveler interface {
	Level() Level
}
