package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger 创建结构化日志记录器
func NewLogger(level, format, output, filePath string) (*zap.Logger, error) {
	var cfg zap.Config

	if format == "console" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
	}

	// 设置日志级别
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
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	// 设置输出
	if output == "file" && filePath != "" {
		cfg.OutputPaths = []string{filePath}
		cfg.ErrorOutputPaths = []string{filePath}
	} else {
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
	}

	// 时间格式
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// Sync 同步日志缓冲区
func Sync(logger *zap.Logger) {
	_ = logger.Sync()
}
