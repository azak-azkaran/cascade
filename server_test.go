package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"testing"
	"time"
)

func Test_createServer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	testServer := createServer(DIRECT.run(true), "localhost", "8082")
	go func() {
		err := testServer.ListenAndServe()
		if err == nil {
			t.Error("Error while starting server", err)
		}
	}()

	time.Sleep(1 * time.Second)
	_, err := client("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over cascade")
	}

	go func() {
		err = shutdown(1 * time.Second)
		if err != nil {
			t.Error("Error while shuting down", err)
		}
	}()

	_, err = client("http://localhost:8082", "https://www.google.de")
	if err == nil {
		t.Error("Error while client request over cascade")
	}

	testServer = createServer(DIRECT.run(true), "localhost", "8082")
	go func() {
		err := testServer.ListenAndServe()
		if err == nil {
			t.Error("Error while starting server", err)
		}
	}()
}
