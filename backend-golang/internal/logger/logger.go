package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"react-example/backend-golang/config"
	"react-example/backend-golang/internal/ctxutil"
)

var Log *zap.Logger
var LogWriter io.Writer

// InitLogger initializes the global Zap logger with lumberjack rotation
func InitLogger() {
	appName := config.AppConfig.AppName
	if appName == "" {
		appName = "iam-governance"
	}

	var logFile *lumberjack.Logger
	// Create logs folder if log to file is enabled
	if config.AppConfig.LogToFile {
		if err := os.MkdirAll("logs", 0755); err != nil {
			fmt.Printf("Failed to create log directory logs: %v\n", err)
		}
		logFile = &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/%s-%s.log", appName, time.Now().Format("2006-01-02")),
			MaxAge:     config.AppConfig.LogRetentionDays, // default 21 days
			Compress:   true,
		}
	}

	// Setup writers for custom integrations (like Fiber)
	var writers []io.Writer
	if config.AppConfig.LogToStdout {
		writers = append(writers, os.Stdout)
	}
	if config.AppConfig.LogToFile && logFile != nil {
		writers = append(writers, logFile)
	}
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}
	LogWriter = io.MultiWriter(writers...)

	// Determine log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.AppConfig.LogLevel)); err != nil {
		level = zapcore.InfoLevel
	}

	// Encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "@timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	var cores []zapcore.Core

	// Configure stdout core
	if config.AppConfig.LogToStdout {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	// Configure file core with lumberjack
	if config.AppConfig.LogToFile && logFile != nil {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(logFile), level))
	}

	// If no core is enabled, just fallback to stdout
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	Log = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(Log)
}

// For returns a contextual logger populated with request ID if available
func For(ctx context.Context) *zap.Logger {
	if Log == nil {
		return zap.L()
	}
	reqID := ctxutil.GetRequestID(ctx)
	if reqID != "" {
		return Log.With(zap.String("request_id", reqID))
	}
	return Log
}
