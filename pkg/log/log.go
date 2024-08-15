package log

import (
	"gin-plus/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

type Level = zapcore.Level

const (
	LevelDebug  = zapcore.DebugLevel
	LevelInfo   = zapcore.InfoLevel
	LevelWarn   = zapcore.WarnLevel
	LevelError  = zapcore.ErrorLevel
	LevelDPanic = zapcore.DPanicLevel
	LevelPanic  = zapcore.PanicLevel
	LevelFatal  = zapcore.FatalLevel
)

type Logger struct {
	l *zap.Logger
}

func Init() {
	//1.Encoder:编码器(如何写入日志)
	zap.NewProductionEncoderConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	logSavePath := config.Conf.Log.SaveDir + config.Conf.Log.Filename + config.Conf.Log.Ext
	//2.writeSyncer:指定将日志写入位置
	writer := &lumberjack.Logger{
		Filename:   logSavePath,
		MaxSize:    config.Conf.Log.MaxSize,    //单位:兆(Mb)
		MaxAge:     config.Conf.Log.MaxAge,     //单位:天
		MaxBackups: config.Conf.Log.MaxBackups, //日志备份数量
		LocalTime:  true,
		Compress:   false,
	}
	writeSyncer := zapcore.AddSync(writer)

	//3.Log Level:记录日志级别
	level := new(zapcore.Level)
	if err := level.UnmarshalText([]byte(config.Conf.Log.Level)); err != nil {
		log.Fatalf("unmarshal log level error: %v", err)
	}

	var core zapcore.Core
	if config.Conf.App.Mode == config.ModeDebug {
		//开发模式下，日志同时输出到控制台和写入文件
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	//替换zap库中全局的logger,因为默认的全局logger什么都不干
	//替换完成后，我们在其他地方就可以通过zap.L()或者zap.S()获取全局logger去记录日志了
	//zap.L()和zap.S()方法可以获取全局的Logger和SugaredLogger
	zap.ReplaceGlobals(logger)
	std = &Logger{l: logger}
}

func New(writer io.Writer, level Level) *Logger {
	if writer == nil {
		writer = os.Stdout
	}
	//1.Encoder:编码器(如何写入日志)
	zap.NewProductionEncoderConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	//2.writeSyncer:指定将日志写入位置
	writeSyncer := zapcore.AddSync(writer)

	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	//替换zap库中全局的logger,因为默认的全局logger什么都不干
	//替换完成后，我们在其他地方就可以通过zap.L()或者zap.S()获取全局logger去记录日志了
	//zap.L()和zap.S()方法可以获取全局的Logger和SugaredLogger
	zap.ReplaceGlobals(logger)
	return &Logger{l: logger}
}

var std = New(os.Stdout, LevelDebug)

func Default() *Logger { return std }

func ReplaceGlobals(logger *Logger) {
	std = logger
	zap.ReplaceGlobals(logger.l)
}

type Field = zapcore.Field

func (this *Logger) Debug(msg string, fields ...Field) {
	this.l.Debug(msg, fields...)
}

func (this *Logger) Info(msg string, fields ...Field) {
	std.l.Info(msg, fields...)
}

func (this *Logger) Warn(msg string, fields ...Field) {
	std.l.Warn(msg, fields...)
}

func (this *Logger) Error(msg string, fields ...Field) {
	std.l.Error(msg, fields...)
}

func (this *Logger) DPanic(msg string, fields ...Field) {
	std.l.DPanic(msg, fields...)
}

func (this *Logger) Panic(msg string, fields ...Field) {
	std.l.Panic(msg, fields...)
}

func (this *Logger) Fatal(msg string, fields ...Field) {
	std.l.Fatal(msg, fields...)
}

func (this *Logger) Sync() error {
	return this.l.Sync()
}

func Debug(msg string, fields ...Field) { std.Debug(msg, fields...) }

func Info(msg string, fields ...Field) { std.Info(msg, fields...) }

func Warn(msg string, fields ...Field) { std.Warn(msg, fields...) }

func Error(msg string, fields ...Field) { std.Error(msg, fields...) }

func DPanic(msg string, fields ...Field) { std.DPanic(msg, fields...) }

func Panic(msg string, fields ...Field) { std.Panic(msg, fields...) }

func Fatal(msg string, fields ...Field) { std.Fatal(msg, fields...) }

func Sync() error { return std.Sync() }
