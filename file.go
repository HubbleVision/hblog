package hblog

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileConfig holds configuration specific to file-based logging
type FileConfig struct {
	// Common config options
	Config
	// Filename is the file to write logs to
	Filename string `toml:"filename" json:"filename"`
	// MaxSize is the maximum size in megabytes of the log file before it gets rotated
	MaxSize int `toml:"max_size" json:"max_size"`
	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int `toml:"max_backups" json:"max_backups"`
	// MaxAge is the maximum number of days to retain old log files
	MaxAge int `toml:"max_age" json:"max_age"`
	// Compress determines if the rotated log files should be compressed
	Compress bool `toml:"compress" json:"compress"`
}

// DefaultFileConfig returns a default configuration for file-based logging
func DefaultFileConfig() *FileConfig {
	return &FileConfig{
		Config:     *DefaultConfig(),
		Filename:   "logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	}
}

// NewFile creates a new logger that writes to a file with rotation
func NewFile(cfg *FileConfig) (*zap.Logger, error) {
	if cfg == nil {
		cfg = DefaultFileConfig()
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Parse log level
	level := ParseLogLevel(cfg.LogLevel)

	// Configure encoder
	encoder := NewJSONEncoder(&cfg.Config)

	// Set up log rotation
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,    // MB
		MaxBackups: cfg.MaxBackups, // number of backups
		MaxAge:     cfg.MaxAge,     // days
		Compress:   cfg.Compress,   // disabled by default
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(lumberJackLogger),
		level,
	)

	// Create and return logger
	return zap.New(core, StandardOptions(&cfg.Config)...), nil
}

// NewFileDefault creates a new file logger with default configuration
func NewFileDefault() (*zap.Logger, error) {
	return NewFile(DefaultFileConfig())
}
