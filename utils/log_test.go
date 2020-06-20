package utils

import (
	"fmt"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {

	message := "message"
	Init()

	Sugar.Info(message)
	Sugar.Error(message)

	//Info.Println(message)
	//s := infobuffer.String()

	//if !strings.Contains(s, "INFO") {
	//	t.Errorf("INFO Buffer does not contain INFO keyword\n%s", s)
	//}
	//if !strings.Contains(s, message) {
	//	t.Errorf("INFO Buffer does not contain message keyword\n%s", s)
	//}
	//if strings.Contains(s, "ERROR") {
	//	t.Errorf("INFO Buffer does contain ERROR keyword\n%s", s)
	//}

	//Error.Println(message)
	//s = errorbuffer.String()

	//if !strings.Contains(s, "ERROR") {
	//	t.Errorf("ERROR Buffer does not contain ERROR keyword\n%s", s)
	//}
	//if !strings.Contains(s, message) {
	//	t.Errorf("ERROR Buffer does not contain message keyword\n%s", s)
	//}
	//if strings.Contains(s, "INFO") {
	//	t.Errorf("ERROR Buffer does contain INFO keyword\n%s", s)
	//}

	//sugar.Infof("test: %s", "blub")
}

func TestSetLogPath(t *testing.T) {
	fmt.Println("Running: TestSetLogPath")
	return

	/*
		message := "message"
		path := "testInfoBuffer"
			buffer := SetLogPath(path)
			LogFile = buffer

			Info.Println(message)
			Warning.Println(message)
			Error.Println(message)

			err := LogFile.Close()
			if err != nil {
				t.Errorf("%s could not be closed: %s", path, err)
			}
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf("error opening file: %v", err)
			}

			m := string(dat)
			if !strings.Contains(m, message) {
				t.Error("File does not contain message")
			}

			if !strings.Contains(m, "INFO") {
				t.Error("File does not contain INFO message")
			}
			if !strings.Contains(m, "WARNING") {
				t.Error("File does not contain WARNING message: ", m)
			}
			if !strings.Contains(m, "ERROR") {
				t.Error("File does not contain ERROR message")
			}

			err = os.Remove(path)
			if err != nil {
				t.Errorf("%s could not be deleted", path)
			}
	*/
}

func TestZapLogger(t *testing.T) {
	fmt.Println("Running: TestZapLogger")

	atom := zap.NewAtomicLevel()
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.Encoding = "console"
	config.Level = atom

	logger, _ := config.Build()

	url := "www.google.de"
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	atom.SetLevel(zap.WarnLevel)
	sugar.Info("Now logs should be colored")
	sugar.Warnw("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Errorf("Failed to fetch URL: %s", url)
}
