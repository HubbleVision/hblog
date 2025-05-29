package hblog

// Config holds common configuration options for all logger types
type Config struct {
	// LogLevel is the minimum log level that will be logged
	LogLevel string `toml:"log_level" json:"log_level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally
	Development bool `toml:"development" json:"development"`
}

// DefaultConfig returns a default configuration with sensible values
func DefaultConfig() *Config {
	return &Config{
		LogLevel:    "info",
		Development: false,
	}
}
