package main

import (
	"context"
	"github.com/azak-azkaran/cascade/utils"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestCreateServer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	testServer := CreateServer(DIRECT.Run(true), "localhost", "8082")
	go func() {
		err := testServer.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error("Error while running server", err)
		}

	}()

	time.Sleep(1 * time.Second)
	resp, err := utils.GetResponse("http://localhost:8082", "http://www.google.de")
	if err != nil {
		t.Error("Error while client request over proxy server", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Error while client request over proxy server", resp.StatusCode)
	}

	err = testServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down", err)
	}

	time.Sleep(1 * time.Second)
	resp, err = utils.GetResponse("http://localhost:8082", "http://www.google.de")
	if err == nil {
		t.Error("No error received on shutdown server", err)
	}
	if resp != nil {
		t.Error("No error received on shutdown server", resp.StatusCode)
	}
}

func TestRunServer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	testServer := CreateServer(DIRECT.Run(true), "localhost", "8082")
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

	err = testServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down server, ", err)
	}

	running = false
	cascade = false
	CONFIG.CascadeFunction = func() {
		cascade = true
		testServer = CreateServer(DIRECT.Run(true), "localhost", "8082")
		CURRENT_SERVER = testServer
	}

	CONFIG.DirectFunction = CONFIG.CascadeFunction
	CURRENT_SERVER = nil
	RunServer()

	time.Sleep(1 * time.Second)
	if !running  {
		t.Error("Server was not started")
	}

	if !cascade {
		t.Error("Cascade function was not called")
	}
	err = testServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down server, ", err)
	}
	running = false
	CURRENT_SERVER = nil
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
	if CURRENT_SERVER != nil {
		t.Error("Server was not removed")
	}
}
