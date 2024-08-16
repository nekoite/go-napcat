package gonapcat

import (
	"sync"

	"github.com/nekoite/go-napcat/config"
)

var initOnce sync.Once

func Init(logCfg *config.LogConfig) {
	initOnce.Do(func() {
		initLogger(logCfg)
	})
}

func Finalize() {
	defer finalizeLogger()
}
