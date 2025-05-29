package hblog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Helper functions for creating common encoder configurations

// NewJSONEncoder creates a JSON encoder with the given configuration
func NewJSONEncoder(cfg *Config) zapcore.Encoder {
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use a more human-friendly encoding in development mode
	if cfg != nil && cfg.Development {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(encCfg)
	}

	return zapcore.NewJSONEncoder(encCfg)
}

// ParseLogLevel parses a log level string into a zapcore.Level
func ParseLogLevel(levelStr string) zapcore.Level {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return zapcore.InfoLevel // default to info level on error
	}
	return level
}

// StandardOptions returns common logger options
func StandardOptions(cfg *Config) []zap.Option {
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if cfg != nil && cfg.Development {
		options = append(options, zap.Development())
	}

	return options
}
