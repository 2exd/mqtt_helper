package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type logger = zap.SugaredLogger

var Logger *logger

func InitLogger() {
	hook := lumberjack.Logger{
		Filename:   "./output.log", // ⽇志⽂件路径
		MaxSize:    1,              // megabytes
		MaxBackups: 3,              // 最多保留 3 个备份
		MaxAge:     7,              // days
		Compress:   false,          // 是否压缩 disabled by default
	}
	var level zapcore.Level

	logLevel := "info"

	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeTime = customTimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		// 打印到控制台和文件
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		// 打印到文件
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)),
		level,
	)

	// 附带堆栈信息
	// l := zap.New(core, zap.AddCaller())
	// 不附带调用者信息
	l := zap.New(core, zap.AddCallerSkip(0))
	Logger = l.Sugar()
}

// 自定义时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
