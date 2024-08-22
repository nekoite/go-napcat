package utils

import (
	"time"
)

func TimedFunc[T any](f func() T, onEnd func(t time.Duration)) T {
	t := time.Now()
	r := f()
	onEnd(time.Since(t))
	return r
}

func TimedAction(f func(), onEnd func(t time.Duration)) {
	t := time.Now()
	f()
	onEnd(time.Since(t))
}
