package hblog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// ConsoleConfig holds configuration specific to console-based logging
type ConsoleConfig struct {
	// Common config options
	Config
	// UseColor determines if log levels should be color-coded
	UseColor bool `toml:"use_color" json:"use_color"`
	// UseJSON determines if console output should be in JSON format
	UseJSON bool `toml:"use_json" json:"use_json"`
}

// DefaultConsoleConfig returns a default configuration for console-based logging
func DefaultConsoleConfig() *ConsoleConfig {
	return &ConsoleConfig{
		Config:   *DefaultConfig(),
		UseColor: true,
		UseJSON:  false,
	}
}

// NewConsole creates a new logger that writes to the console (stdout)
func NewConsole(cfg *ConsoleConfig) (*zap.Logger, error) {
	if cfg == nil {
		cfg = DefaultConsoleConfig()
	}

	// Parse log level
	level := ParseLogLevel(cfg.LogLevel)

	var encoder zapcore.Encoder
	if cfg.UseJSON {
		encoder = NewJSONEncoder(&cfg.Config)
	} else {
		// Create a more human-friendly encoder for console output
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		if cfg.UseColor {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		}

		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create and return logger
	return zap.New(core, StandardOptions(&cfg.Config)...), nil
}

// NewConsoleDefault creates a new console logger with default configuration
func NewConsoleDefault() (*zap.Logger, error) {
	return NewConsole(DefaultConsoleConfig())
}
