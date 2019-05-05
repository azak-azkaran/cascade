package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"testing"
	"time"
)

func TestCreateServer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	testServer := CreateServer(DIRECT.Run(true), "localhost", "8082")
	go func() {
		err := testServer.ListenAndServe()
		if err == nil {
			t.Error("Error while starting server", err)
		}
	}()

	time.Sleep(1 * time.Second)
	resp, err := utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over proxy server", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Error while client request over proxy server", resp.StatusCode)
	}

	go func() {
		err = shutdown(1 * time.Second)
		if err != nil {
			t.Error("Error while shuting down", err)
		}
	}()

	resp, err = utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err == nil {
		t.Error("No error received on shutdown server", err)
	}
	if resp != nil {
		t.Error("No error received on shutdown server", resp.StatusCode)
	}

	testServer = CreateServer(DIRECT.Run(true), "localhost", "8082")
	go func() {
		err := testServer.ListenAndServe()
		if err == nil {
			t.Error("Error while starting server", err)
		}
	}()
}

func TestRunServer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CreateServer(DIRECT.Run(true), "localhost", "8082")
	if running {
		t.Error("Server already running")
	}

	RunServer()
	time.Sleep(1 * time.Second)
	if !running {
		t.Error("Server was not started")
	}

	resp, err := utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over proxy server", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Error while client request over proxy server", resp.StatusCode)
	}
}

func TestShutdownCurrentServer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CreateServer(DIRECT.Run(true), "localhost", "8082")
	if running {
		t.Error("Server already running")
	}

	RunServer()
	time.Sleep(1 * time.Second)
	if !running {
		t.Error("Server was not started")
	}

	ShutdownCurrentServer()
	time.Sleep(1 * time.Second)
	if running {
		t.Error("Server was not shutdown")
	}
}
