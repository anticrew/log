package log

import "errors"

var (
	ErrUnknownLevel     = errors.New("unknown level")
	ErrUnknownPrintMode = errors.New("unknown print mode")
)

var (
	ErrUnknownKind = errors.New("unknown kind")
)
