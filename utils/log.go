package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Atom      = zap.NewAtomicLevel()
	LogConfig = zap.NewProductionConfig()
	logger    *zap.Logger
	Sugar     *zap.SugaredLogger
)

func Init() {
	var err error
	LogConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	LogConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	LogConfig.Encoding = "console"
	LogConfig.Level = Atom
	logger, err = LogConfig.Build()
	if err != nil {
		fmt.Println("Error building logger:", err)
	}
	defer logger.Sync() // flushes buffer, if any
	Sugar = logger.Sugar()

	gin.SetMode(gin.ReleaseMode)
	//gin.DisableConsoleColor()
	//gin.DefaultErrorWriter = errorHandle
	//std := zap.NewStdLog(logger)
	//var writer io.Writer
	//gin.DefaultWriter = writer
	//std.SetOutput(writer)
}

func GetLogger() *log.Logger {
	return zap.NewStdLog(logger)
}

func EnableError() {
	Atom.SetLevel(zap.ErrorLevel)
}

func EnableWarning() {
	Atom.SetLevel(zap.WarnLevel)
}

func EnableInfo() {
	Atom.SetLevel(zap.InfoLevel)
}

func EnableDebug() {
	Atom.SetLevel(zap.DebugLevel)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var DefaultLogFormatter = func(param gin.LogFormatterParams) string {
	var methodColor, resetColor string
	//if param.IsOutputColor() {
	methodColor = param.MethodColor()
	resetColor = param.ResetColor()
	//}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("GIN: %s %-7s %s| %13v | %15s | %s\n%s",
		//param.TimeStamp.Format("2006/01/02 | 15:04:05"),
		methodColor, param.Method, resetColor,
		param.Latency,
		param.Request.URL.Host,
		param.Path,
		param.ErrorMessage,
	)
}

//func SetLogPath(path string) *os.File {
//	buffer, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	zap.Open(path)
//	if err != nil {
//		Error.Println("Error while opining Log file:", err)
//		return nil
//	}
//	Init(buffer, buffer, buffer)
//	return buffer
//}

//func Close() {
//	if LogFile != nil {
//		err := LogFile.Close()
//		if err != nil {
//			Error.Println("Error while closing LogFile Pointer: ", err)
//		}
//	}
//}
