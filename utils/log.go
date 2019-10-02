package utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	Info          *log.Logger = disableLogger()
	Warning       *log.Logger = disableLogger()
	Error         *log.Logger = disableLogger()
	Discard       *log.Logger = disableLogger()
	infoWriter    io.Writer   = os.Stdout
	warningWriter io.Writer   = os.Stdout
	errorWriter   io.Writer   = os.Stderr
)

func Init(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	infoWriter = infoHandle
	EnableInfo()

	warningWriter = warningWriter
	EnableWarning()

	errorWriter = errorHandle
	EnableError()
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
