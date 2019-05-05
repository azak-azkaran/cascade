package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"net/http"
	"os"
	"testing"
	"time"
)

func Test_run(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	directProxy := DIRECT.run(true)
	go http.ListenAndServe("localhost:8082", directProxy)

	time.Sleep(1 * time.Second)
	utils.Info.Println("waiting for running")

	dump, err := client("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over cascade")
	}

	if len(dump) == 0 {
		t.Error("Error while client Request, dump is empty")
	}

	dump, err = client("http://localhost:8082", "http://www.google.de")
	if err != nil {
		t.Error("Error while client request over cascade")
	}

	if len(dump) == 0 {
		t.Error("Error while client Request, dump is empty")
	}
}
