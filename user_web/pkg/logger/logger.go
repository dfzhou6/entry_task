package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"user_web/pkg/config"
)

var Logger *zap.Logger

// SetupLogger 初始化日志库
func SetupLogger() {
	var coreArr []zapcore.Core

	// 获取编码器
	var encoderConfig zapcore.EncoderConfig
	if config.GetBool("APP_DEBUG") {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// 日志级别
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	// info文件writeSyncer
	infoFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/info.log", // 日志文件存放目录，如果文件夹不存在会自动创建
		MaxSize:    2,                 // 文件大小限制,单位MB
		MaxBackups: 100,               // 最大保留日志文件数量
		MaxAge:     30,                // 日志文件保留天数
		Compress:   false,             // 是否压缩处理
	})
	infoFileCore := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(infoFileWriteSyncer, zapcore.AddSync(os.Stdout)), lowPriority)
	// error文件writeSyncer
	errorFileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/error.log", // 日志文件存放目录
		MaxSize:    1,                  // 文件大小限制,单位MB
		MaxBackups: 5,                  // 最大保留日志文件数量
		MaxAge:     30,                 // 日志文件保留天数
		Compress:   false,              // 是否压缩处理
	})
	errorFileCore := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(errorFileWriteSyncer, zapcore.AddSync(os.Stdout)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)
	Logger = zap.New(zapcore.NewTee(coreArr...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.WarnLevel))
}

func Debug(module string, msg ...interface{}) {
	Logger.Sugar().Debug("module=", module, msg)
}

func Info(module string, msg ...interface{}) {
	Logger.Sugar().Info("module=", module, msg)
}

func Warn(module string, msg ...interface{}) {
	Logger.Sugar().Warn("module=", module, msg)
}

func Error(module string, msg ...interface{}) {
	Logger.Sugar().Error("module=", module, msg)
}
