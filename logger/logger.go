package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(level string) error {
    var zapLevel zapcore.Level
    switch level {
    case "debug":
        zapLevel = zapcore.DebugLevel
    case "info":
        zapLevel = zapcore.InfoLevel
    case "warn":
        zapLevel = zapcore.WarnLevel
    case "error":
        zapLevel = zapcore.ErrorLevel
    default:
        zapLevel = zapcore.InfoLevel
    }

    config := zap.Config{
        Level:       zap.NewAtomicLevelAt(zapLevel),
        Development: false,
        Encoding:    "json",
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "timestamp",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "message",
            StacktraceKey:  "stacktrace",
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeLevel:    zapcore.LowercaseLevelEncoder,
            EncodeTime:     zapcore.ISO8601TimeEncoder,
            EncodeDuration: zapcore.SecondsDurationEncoder,
            EncodeCaller:   zapcore.ShortCallerEncoder,
        },
        OutputPaths:      []string{"stdout"},
        ErrorOutputPaths: []string{"stderr"},
    }

    var err error
    Log, err = config.Build()
    if err != nil {
        return err
    }

    return nil
}

func Sync() {
    if Log != nil {
        Log.Sync()
    }
}

func Info(msg string, fields ...zap.Field) {
    if Log != nil {
        Log.Info(msg, fields...)
    }
}

func Error(msg string, fields ...zap.Field) {
    if Log != nil {
        Log.Error(msg, fields...)
    }
}

func Debug(msg string, fields ...zap.Field) {
    if Log != nil {
        Log.Debug(msg, fields...)
    }
}

func Warn(msg string, fields ...zap.Field) {
    if Log != nil {
        Log.Warn(msg, fields...)
    }
}

func Fatal(msg string, fields ...zap.Field) {
    if Log != nil {
        Log.Fatal(msg, fields...)
    }
}
