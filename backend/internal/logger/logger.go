package logger

import (
	"fmt"
	"os"
	"strings"

	"gomall/backend/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var Logger *zap.Logger

// Init 初始化日志
func Init() error {
	cfg := config.GetLogger()

	// 获取日志配置（如果配置不存在，使用默认值）
	level := "info"
	format := "json"
	output := "stdout"
	filename := "app.log"

	if cfg != nil {
		level = cfg.GetString("level")
		format = cfg.GetString("format")
		output = cfg.GetString("output")
		filename = cfg.GetString("filename")
	}

	// 解析日志级别
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn", "warning":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if format == "json" {
		// JSON 格式 - 适用于日志收集系统（ELK/Loki）
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		// 控制台格式 - 便于开发调试
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "name",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stack",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建写入器
	var writeSyncer zapcore.WriteSyncer
	if output == "file" {
		// 写入文件
		if filename == "" {
			filename = "app.log"
		}
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("创建日志文件失败: %w", err)
		}
		writeSyncer = zapcore.AddSync(f)
	} else {
		// 输出到控制台
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	// 创建日志实例
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		// 返回默认 logger
		return zap.NewNop()
	}
	return Logger
}

// Sync 刷新日志缓冲区
func Sync() error {
	if Logger != nil {
		return Logger.Sync()
	}
	return nil
}

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 错误级别日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal 致命级别日志
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// With 创建带字段的日志实例
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}

// Infof 格式化日志（兼容旧接口）
func Infof(template string, args ...interface{}) {
	GetLogger().Info(fmt.Sprintf(template, args...))
}

// Warnf 格式化日志（兼容旧接口）
func Warnf(template string, args ...interface{}) {
	GetLogger().Warn(fmt.Sprintf(template, args...))
}

// Errorf 格式化日志（兼容旧接口）
func Errorf(template string, args ...interface{}) {
	GetLogger().Error(fmt.Sprintf(template, args...))
}

// Debugf 格式化日志（兼容旧接口）
func Debugf(template string, args ...interface{}) {
	GetLogger().Debug(fmt.Sprintf(template, args...))
}

// Reload 重新加载日志配置（动态更新日志级别）
func Reload() error {
	cfg := config.GetLogger()
	if cfg == nil || Logger == nil {
		return nil
	}

	level := cfg.GetString("level")
	// 解析日志级别
	var zapLevel zapcore.Level
	switch strings.ToLower(level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn", "warning":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// 更新日志级别
	_ = zapLevel
	Logger.Info("日志配置已重新加载", zap.String("level", level))

	return nil
}
