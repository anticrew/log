package main

import (
	"errors"
	"os"

	"github.com/anticrew/log/engine"
	"github.com/anticrew/log/middleware/interpolation"

	"github.com/anticrew/log"
)

func main() {
	base := log.New(engine.NewEngine(os.Stdout, log.PrintModePretty), log.LevelTrace, interpolation.New())

	log.SetDefault(base.With(log.String("instance", "default")))
	l := base.With(log.String("instance", "local"))

	log.Trace("trace message")
	l.Trace("trace message")

	log.Debug("debug message")
	l.Debug("debug message")

	log.Info("info message")
	l.Debug("info message")

	log.Warn("warning message")
	l.Warn("warning message")

	log.Debug(`debug message with "{{ key }}" and attr`, log.String("key", "interpolation"))
	l.Debug(`debug message with "{{ key }}" and attr`, log.String("key", "interpolation"))

	log.Debug("debug message with attr only", log.String("key", "interpolation"))
	l.Debug("debug message with attr only", log.String("key", "interpolation"))

	log.Warn(`warn message with interpolation from "{{ instance }}" instance`)
	l.Warn(`warn message with interpolation from "{{ instance }}" instance`)

	err := errors.New("not found")
	log.Error(err, "oh, error?")
	l.Error(err, "oh, error?")

	log.Fatal(err, "oh, fail!")

	panic("why im called?") // never called because of xlog.Fatal
}
