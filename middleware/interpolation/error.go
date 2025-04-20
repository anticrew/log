package interpolation

import "errors"

var (
	ErrIndirectOpenKey  = errors.New("indirect open key")
	ErrIndirectCloseKey = errors.New("indirect close key")
	ErrEmptyKey         = errors.New("empty key")
	ErrUnclosedKey      = errors.New("unclosed key")
)
