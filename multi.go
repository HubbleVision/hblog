package hblog

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MultiConfig holds configuration for a logger that outputs to multiple destinations
type MultiConfig struct {
	// Common config options
	Config
	// ConsoleEnabled determines if console logging is enabled
	ConsoleEnabled bool `toml:"console_enabled" json:"console_enabled"`
	// ConsoleJSON determines if console output should be in JSON format
	ConsoleJSON bool `toml:"console_json" json:"console_json"`
	// ConsoleColor determines if console output should use colors
	ConsoleColor bool `toml:"console_color" json:"console_color"`
	// FileEnabled determines if file logging is enabled
	FileEnabled bool `toml:"file_enabled" json:"file_enabled"`
	// Filename is the file to write logs to
	Filename string `toml:"filename" json:"filename"`
}

// DefaultMultiConfig returns a default configuration for multi-destination logging
func DefaultMultiConfig() *MultiConfig {
	return &MultiConfig{
		Config:         *DefaultConfig(),
		ConsoleEnabled: true,
		ConsoleJSON:    false,
		ConsoleColor:   true,
		FileEnabled:    true,
		Filename:       "logs/app.log",
	}
}

// NewMulti creates a new logger that writes to multiple destinations (console and file)
func NewMulti(cfg *MultiConfig) (*zap.Logger, error) {
	if cfg == nil {
		cfg = DefaultMultiConfig()
	}

	// Parse log level
	level := ParseLogLevel(cfg.LogLevel)

	// Prepare cores slice for each output
	var cores []zapcore.Core

	// Add console core if enabled
	if cfg.ConsoleEnabled {
		var consoleEncoder zapcore.Encoder
		if cfg.ConsoleJSON {
			consoleEncoder = NewJSONEncoder(&cfg.Config)
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

			if cfg.ConsoleColor {
				encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			} else {
				encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
			}

			consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// Add file core if enabled
	if cfg.FileEnabled {
		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		fileEncoder := NewJSONEncoder(&cfg.Config)
		// Open the file
		file, err := os.OpenFile(cfg.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		fileCore := zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(file),
			level,
		)
		cores = append(cores, fileCore)
	}

	// Check if we have at least one core
	if len(cores) == 0 {
		return nil, fmt.Errorf("no logging destination enabled")
	}

	// Create and return logger with multiple cores
	core := zapcore.NewTee(cores...)
	return zap.New(core, StandardOptions(&cfg.Config)...), nil
}

// NewMultiDefault creates a new multi-destination logger with default configuration
func NewMultiDefault() (*zap.Logger, error) {
	return NewMulti(DefaultMultiConfig())
}
