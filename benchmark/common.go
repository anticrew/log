package benchmark

import "github.com/anticrew/log"

type Entity log.Record

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
