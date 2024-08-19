package gonapcat

import (
	"github.com/nekoite/go-napcat/config"
	"github.com/nekoite/go-napcat/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func initLogger(cfg *config.LogConfig) {
	var zapConfig zap.Config
	if cfg.Debug {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	err := zapConfig.Level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(err)
	}
	zapConfig.OutputPaths = make([]string, 0, 1)
	if len(cfg.Paths) > 0 {
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, cfg.Paths...)
	}
	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	l, err := zapConfig.Build(zap.AddCaller())
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l.Named(consts.AppName))
	logger = l
}

func finalizeLogger() {
	if logger != nil {
		logger.Sync()
	}
}
