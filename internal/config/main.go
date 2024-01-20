package config

import (
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  "time"
)

func initLogger(option string) *zap.Logger {
  config := zap.NewDevelopmentConfig()
  config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

  if option == "prod" {
    config = zap.NewProductionConfig()
    config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
  }

  config.EncoderConfig.TimeKey = "timestamp"
  config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.UTC().Format(time.RFC3339))
  }
  logger, err := config.Build()
  if err != nil {
    panic("cannot initialize logger-create")
  }
  return logger
}
