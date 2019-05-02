package main

import (
	"bytes"
	"github.com/azak-azkaran/proxy-go/utils"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	var infobuffer bytes.Buffer
	utils.Init(&infobuffer, os.Stdout, os.Stderr)
	if !created {
		go run(":8889")
	}

	for !running {
		time.Sleep(1 * time.Second)
		utils.Info.Println("waiting for running")
	}
}

func Test_handleHTTP(t *testing.T) {
	if !running {
		Test_run(t)
	}
	var infobuffer bytes.Buffer
	utils.Init(&infobuffer, os.Stdout, os.Stderr)

	dump, _ := client("http://localhost:8889", "http://google.de")

	if len(dump) == 0 {
		t.Error("No dump received")
	}

	logMessages := infobuffer.String()
	if !strings.Contains(logMessages, "GET") {
		t.Error("No http request received")
	}

	client("http://localhost:8889", "http://localhost:12313")
	if strings.Contains(infobuffer.String(), "503") {
		t.Error("Server available but should not", infobuffer.String())
	}

}

func Test_handleTunneling(t *testing.T) {
	if !running {
		Test_run(t)
	}
	var infobuffer bytes.Buffer
	utils.Init(&infobuffer, os.Stdout, os.Stderr)

	dump, _ := client("http://localhost:8889", "https://www.google.de")

	if len(dump) == 0 {
		t.Error("No dump recieved")
	}

	logMessages := infobuffer.String()
	if strings.Contains(logMessages, "GET") {
		t.Error("HTTP Request received instead of HTTPS")
	}

	if !strings.Contains(logMessages, "CONNECT") {
		t.Error("Did not receive a HTTPS Request")
	}
	client("http://localhost:8889", "https://localhost:12313")
	if strings.Contains(infobuffer.String(), "503") {
		t.Error("Server available but should not", infobuffer.String())
	}
}

func Test_shutdown(t *testing.T) {
	if !running {
		Test_run(t)
	}
	var errorbuffer bytes.Buffer
	utils.Init(os.Stdout, os.Stdout, &errorbuffer)
	go shutdown(5 * time.Second)
	time.Sleep(5 * time.Second)
	if running {
		t.Error("Was not shutdown in time")
	}

	str := errorbuffer.String()
	if !strings.Contains(str, "Server closed") {
		t.Error("Server was not closed\n", str)
	}
	err := shutdown(1 * time.Nanosecond)
	if err == nil {
		t.Error("Could shutdown server twice")
	}
}
