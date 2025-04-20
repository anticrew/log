package caller

import (
	"runtime"

	"github.com/anticrew/go-x/pool"
)

// _pool
// Пул для переиспользования экземпляров Caller
var _pool = pool.NewPool(func() *Caller {
	return &Caller{}
})

// Caller
// Структура, хранящая информацию об источнике вызова
type Caller struct {
	Function string
	File     string
	Line     int

	pcs [1]uintptr
}

// Capture
// Захватывает информацию об источнике вызова, пропуская 2 вызова: вызов самого Capture и runtime.Callers
func Capture(skipCount int) (*Caller, func(), error) {
	// always skip runtime.Callers and stack.Capture
	skipCount += 2

	caller := _pool.Get()

	if framesCount := runtime.Callers(skipCount, caller.pcs[:]); framesCount < 1 {
		return nil, nil, ErrNoFrames
	}

	f, _ := runtime.CallersFrames(caller.pcs[:1]).Next()
	if f.PC == 0 {
		return nil, nil, ErrNoFrames
	}

	caller.File, caller.Line = f.File, f.Line
	caller.Function = f.Function

	return caller, callerDisposer(caller), nil
}

// callerDisposer
// Возвращает функцию для возвращения экземпляра Caller в общий пул
func callerDisposer(c *Caller) func() {
	return func() {
		_pool.Put(c)
	}
}
