package config

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	ShutdownTimeout time.Duration
	LogLevel        slog.Level
}

func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        "127.0.0.1:8080",
		ShutdownTimeout: 5 * time.Second,
		LogLevel:        slog.LevelInfo,
	}
	if value, ok := os.LookupEnv("HTTP_ADDR"); ok {
		_, port, err := net.SplitHostPort(value)
		if err != nil {
			return Config{}, fmt.Errorf("HTTP_ADDR=%q: %w", value, err)
		}
		portInt, err := strconv.Atoi(port)
		if err != nil {
			return Config{}, fmt.Errorf("HTTP_ADDR=%q: %w", value, err)
		}
		if portInt < 1 || portInt > 65535 {
			return Config{}, fmt.Errorf("HTTP_ADDR=%q: invalid port", value)
		}
		cfg.HTTPAddr = value
	}
	if value, ok := os.LookupEnv("SHUTDOWN_TIMEOUT"); ok {
		d, err := time.ParseDuration(value)
		if err != nil {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT=%q: %w", value, err)
		}
		if d <= 0 {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT=%q: invalid timeout", value)
		}
		cfg.ShutdownTimeout = d
	}
	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch strings.ToUpper(value) {
		case "INFO":
			cfg.LogLevel = slog.LevelInfo
		case "DEBUG":
			cfg.LogLevel = slog.LevelDebug
		case "ERROR":
			cfg.LogLevel = slog.LevelError
		case "WARN":
			cfg.LogLevel = slog.LevelWarn
		default:
			return Config{}, fmt.Errorf("LOG_LEVEL=%q: invalid level", value)
		}
	}
	return cfg, nil
}
