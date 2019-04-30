package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {

	message := "message"
	var infobuffer bytes.Buffer
	var warningbuffer bytes.Buffer
	var errorbuffer bytes.Buffer
	Init(&infobuffer, &warningbuffer, &errorbuffer)

	Info.Println(message)
	s := infobuffer.String()

	if !strings.Contains(s, "INFO") {
		t.Errorf("INFO Buffer does not contain INFO keyword\n%s", s)
	}
	if !strings.Contains(s, message) {
		t.Errorf("INFO Buffer does not contain message keyword\n%s", s)
	}
	if strings.Contains(s, "ERROR") {
		t.Errorf("INFO Buffer does contain ERROR keyword\n%s", s)
	}

	Error.Println(message)
	s = errorbuffer.String()

	if !strings.Contains(s, "ERROR") {
		t.Errorf("ERROR Buffer does not contain ERROR keyword\n%s", s)
	}
	if !strings.Contains(s, message) {
		t.Errorf("ERROR Buffer does not contain message keyword\n%s", s)
	}
	if strings.Contains(s, "INFO") {
		t.Errorf("ERROR Buffer does contain INFO keyword\n%s", s)
	}

}
