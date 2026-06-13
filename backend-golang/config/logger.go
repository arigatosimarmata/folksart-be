package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerConfig encapsulates all configuration parameters for the application logger.
// This prevents functions from having too many parameters (go:S107).
type LoggerConfig struct {
	LogPath    string
	MaxAgeDays int
	MaxSizeMB  int
	MaxBackups int
	Level      zapcore.Level
}

// Logger is the globally accessible zap logger instance.
var Logger *zap.Logger

// InitLogger initializes the global zap logger with lumberjack log rotation.
func InitLogger(cfg LoggerConfig) *zap.Logger {
	setLoggerDefaults(&cfg)

	encoder := getJSONEncoder()
	writeSyncer := getWriteSyncer(cfg)

	core := zapcore.NewCore(encoder, writeSyncer, cfg.Level)

	Logger = zap.New(core, zap.AddCaller())
	return Logger
}

// setLoggerDefaults applies default configuration values if not specified.
func setLoggerDefaults(cfg *LoggerConfig) {
	if cfg.LogPath == "" {
		cfg.LogPath = "app.log"
	}
	if cfg.MaxAgeDays == 0 {
		cfg.MaxAgeDays = 5 // Default rotation to 5 days as requested
	}
	if cfg.MaxSizeMB == 0 {
		cfg.MaxSizeMB = 100
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 10
	}
}

// getJSONEncoder returns a configured JSON encoder for Fluentd/Elasticsearch.
func getJSONEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getWriteSyncer returns a multi-write syncer for stdout and rolling file.
func getWriteSyncer(cfg LoggerConfig) zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   true,
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberjackLogger),
	)
}
