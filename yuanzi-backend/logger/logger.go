package logger

import (
	"os"
	"yuanzi-backend/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// Setup 初始化日志
func Setup() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()

	var level zapcore.Level
	switch config.GlobalConfig.Server.RunMode {
	case "debug":
		level = zapcore.DebugLevel
	case "release":
		level = zapcore.InfoLevel
	case "test":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(encoder, writeSyncer, level)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if config.GlobalConfig.Server.RunMode == "debug" {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	// 控制台输出
	consoleSyncer := zapcore.AddSync(os.Stdout)

	// 文件输出
	fileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100, // MB
		MaxBackups: 10,
		MaxAge:     30, // days
		Compress:   true,
	})

	// 错误日志单独输出
	errorFileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	})

	// 使用 MultiWriteSyncer 同时输出到控制台和文件
	return zapcore.NewMultiWriteSyncer(consoleSyncer, fileSyncer, errorFileSyncer)
}

// Debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info 级别日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error 级别日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal 级别日志
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// String 创建字符串字段
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建整数字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Err 创建错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
