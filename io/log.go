package io

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() func() error {
	if Logger != nil {
		Logger.Error("double logger initialization")
		return func() error {
			return nil
		}
	}

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.CallerKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	Logger = zap.Must(config.Build())

	return Logger.Sync
}
