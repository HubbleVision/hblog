package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"github.com/HubbleVision/hblog"
)

func main() {
	exampleFileLogger()
	// Uncomment to try other logger types
	// exampleConsoleLogger()
	// exampleMultiLogger()
	// exampleKafkaLogger()
}

func exampleFileLogger() {
	// Create a file logger with custom configuration
	cfg := hblog.DefaultFileConfig()
	cfg.Filename = "logs/example.log"
	cfg.MaxSize = 10 // 10MB
	cfg.MaxBackups = 5
	cfg.MaxAge = 7 // 7 days

	logger, err := hblog.NewFile(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create file logger: %w", err))
	}
	defer logger.Sync()

	logSampleMessages(logger, "File logger")
}

func exampleConsoleLogger() {
	// Create a console logger with custom configuration
	cfg := hblog.DefaultConsoleConfig()
	cfg.UseColor = true
	cfg.LogLevel = "debug" // Set to debug to see all messages

	logger, err := hblog.NewConsole(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create console logger: %w", err))
	}
	defer logger.Sync()

	logSampleMessages(logger, "Console logger")
}

func exampleMultiLogger() {
	// Create a multi-destination logger with custom configuration
	cfg := hblog.DefaultMultiConfig()
	cfg.Filename = "logs/multi_example.log"
	cfg.ConsoleColor = true
	cfg.LogLevel = "debug"

	logger, err := hblog.NewMulti(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create multi logger: %w", err))
	}
	defer logger.Sync()

	logSampleMessages(logger, "Multi logger")
}

func exampleKafkaLogger() {
	// Example for when Kafka logger is implemented
	cfg := hblog.DefaultKafkaConfig()
	cfg.Brokers = []string{"localhost:9092"}
	cfg.Topic = "application-logs"

	logger, err := hblog.NewKafka(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create kafka logger: %w", err))
	}
	defer logger.Sync()

	logSampleMessages(logger, "Kafka logger")
}

// logSampleMessages logs a set of sample messages at different log levels
func logSampleMessages(logger *zap.Logger, loggerType string) {
	// Log some messages
	logger.Info(loggerType + " initialized successfully")
	logger.Debug("This is a debug message")
	logger.Warn("This is a warning message",
		zap.String("key", "value"),
		zap.Int("count", 42))

	// Example with error
	if err := someFunction(); err != nil {
		logger.Error("Failed to execute someFunction",
			zap.Error(err),
			zap.String("time", time.Now().String()))
	}
}

func someFunction() error {
	return nil
}
