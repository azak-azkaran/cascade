package main

import (
	"context"
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"net/http"
	"os"
	"testing"
	"time"
)

var test_config Yaml = Yaml{LocalPort: "8082", verbose: true}

func TestCreateServer(t *testing.T) {
	fmt.Println("Running: TestCreateServer")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	testServer := CreateServer(test_config) //CASCADE.Run(true, "", "", ""), "localhost", "8082")
	go func() {
		err := testServer.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error("Error while running server", err)
		}

	}()

	time.Sleep(1 * time.Second)
	DirectOverrideChan = true
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
	DirectOverrideChan = false
}

func TestRunServer(t *testing.T) {
	fmt.Println("Running: TestRunServer")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	testServer := CreateServer(test_config)
	if running {
		t.Error("Server already running")
	}

	if CurrentServer == nil {
		t.Error("Server was not created")
	}

	RunServer()
	time.Sleep(1 * time.Second)
	if !running {
		t.Error("Server was not started")
	}
	DirectOverrideChan = true

	resp, err := utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over proxy server", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Error while client request over proxy server", resp)
	}

	err = testServer.Shutdown(context.TODO())
	time.Sleep(1 * time.Second)
	if err != nil {
		t.Error("Error while shutting down server, ", err)
	}
	DirectOverrideChan = false
	CurrentServer = nil
}

func TestShutdownCurrentServer(t *testing.T) {
	fmt.Println("Running: TestShutdownCurrentServer")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CreateServer(test_config)
	if running {
		t.Error("Server already running")
	}

	ShutdownCurrentServer()
	time.Sleep(1 * time.Second)
	if running {
		t.Error("Server was not shutdown")
	}
	if CurrentServer != nil {
		t.Error("Server was not removed")
	}
}

func TestCreateBrokenServer(t *testing.T) {
	fmt.Println("Running: TestCreateBrokenServer")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CreateConfig("8082", "", "", "", "https://www.google.de", 5, "golang.org,youtube.com", "info")

	RunServer()
	time.Sleep(1 * time.Second)
	if !running {
		t.Error("Server was not started")
	}

	if len(Config.proxyRedirectList) != 2 {
		t.Error("Skip for Cascade list was not separated correctly")
	}

	_, err := utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err == nil {
		t.Error("Request over broken proxy was successfull but should not be", err)
	}

	resp, err := utils.GetResponse("http://localhost:8082", "http://golang.org/doc/")
	if err != nil {
		t.Error("Error while client request over proxy server", err)
	}
	if resp != nil && resp.StatusCode != 200 {
		t.Error("Error while client request over proxy server", resp.Status)
	}

	ShutdownCurrentServer()
	time.Sleep(1 * time.Second)
	if running {
		t.Error("Server was not shutdown")
	}
	if CurrentServer != nil {
		t.Error("Server was not removed")
	}
}

func TestRestRequest(t *testing.T) {
	fmt.Println("Running: TestCreateBrokenServer")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	testServer := CreateServer(test_config)

	RunServer()
	time.Sleep(1 * time.Second)
	if !running {
		t.Error("Server was not started")
	}
	DirectOverrideChan = true

	client, err := utils.GetClient("", 2)
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := client.Get("http://localhost:8082/config")
	if err != nil {
		t.Error(err)
		return
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Response should be 200", resp.StatusCode)
	}

	resp, err = utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over proxy server", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Error while client request over proxy server", resp)
	}

	err = testServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down", err)
	}
}
