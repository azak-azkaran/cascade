package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	Info          *log.Logger = disableLogger()
	Warning       *log.Logger = disableLogger()
	Error         *log.Logger = disableLogger()
	Discard       *log.Logger = disableLogger()
	infoWriter    io.Writer   = os.Stdout
	warningWriter io.Writer   = os.Stdout
	errorWriter   io.Writer   = os.Stderr
	// LogFile File for logs if log to file is active
	LogFile *os.File
)

func Init(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	infoWriter = infoHandle
	EnableInfo()

	warningWriter = warningHandle
	EnableWarning()

	errorWriter = errorHandle
	EnableError()

	gin.SetMode(gin.ReleaseMode)
	//gin.DisableConsoleColor()
	gin.DefaultErrorWriter = errorHandle
	//gin.DefaultWriter = infoWriter
}

func EnableError() {
	Error = log.New(errorWriter,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func EnableWarning() {
	Warning = log.New(warningWriter,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func EnableInfo() {
	Info = log.New(infoWriter,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func disableLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}

func DisableInfo() {
	Info = disableLogger()
}

func DisableWarning() {
	Warning = disableLogger()
}

func DisableError() {
	Error = disableLogger()
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
	return fmt.Sprintf("GIN: %v |%s %-7s %s| %13v | %15s | %s\n%s",
		param.TimeStamp.Format("2006/01/02 | 15:04:05"),
		methodColor, param.Method, resetColor,
		param.Latency,
		param.Request.URL.Host,
		param.Path,
		param.ErrorMessage,
	)
}

func SetLogPath(path string) *os.File {
	buffer, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Error.Println("Error while opining Log file:", err)
		return nil
	}
	Init(buffer, buffer, buffer)
	return buffer
}

func Close() {
	if LogFile != nil {
		err := LogFile.Close()
		if err != nil {
			Error.Println("Error while closing LogFile Pointer: ", err)
		}
	}
}
