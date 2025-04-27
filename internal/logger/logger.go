package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// New returns a zap.Logger configured with the given level.
func New(level string) (*zap.Logger, error) {
    var zapLevel zapcore.Level
    switch level {
    case "debug":
        zapLevel = zap.DebugLevel
    case "info":
        zapLevel = zap.InfoLevel
    case "warn":
        zapLevel = zap.WarnLevel
    case "error":
        zapLevel = zap.ErrorLevel
    default:
        zapLevel = zap.InfoLevel
    }
    cfg := zap.NewProductionConfig()
    cfg.Level = zap.NewAtomicLevelAt(zapLevel)
    return cfg.Build()
}
