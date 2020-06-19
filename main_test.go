package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	fmt.Println("Running: TestRun")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	config := Yaml{}
	config.HealthTime = 1
	config.Username = "foo"
	config.Password = "bar"
	config.ProxyURL = "localhost:8082"
	config.LocalPort = "8081"
	config.Log = "info"
	config.CheckAddress = "https://google.de"
	config.ConfigFile = ""
	config.VaultAddr = ""

	go Run(config)

	time.Sleep(1 * time.Second)
	assert.NotNil(t, CurrentServer, "No Server was created")

	assert.True(t, DirectOverrideChan, "Direct Override is not active")

	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cleanup()
	time.Sleep(1 * time.Second)

	assert.False(t, running, "Server is still running")
	assert.Nil(t, CurrentServer)
}

func TestMain(t *testing.T) {
	fmt.Println("Running: TestMain")
	closeChan = false
	config := Yaml{}
	config.HealthTime = 5
	config.Username = "foo"
	config.Password = "bar"
	config.ProxyURL = "localhost:8082"
	config.LocalPort = "8888"
	config.Log = "info"
	config.CheckAddress = "https://google.de"
	config.ConfigFile = ""
	config.VaultAddr = ""
	Config = config

	go main()

	time.Sleep(2 * time.Second)
	utils.Info.Println("calling HTTP")
	resp, err := utils.GetResponse("http://localhost:8888", "http://www.google.de")
	if err != nil {
		t.Error("http test failed: ", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("http test failed: ", resp)
	}

	utils.Info.Println("calling HTTPs")
	resp, err = utils.GetResponse("http://localhost:8888", "https://www.google.de")
	if err != nil {
		t.Error("http test failed: ", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("http test failed: ", resp)
	}

	utils.Info.Println("Closing")
	stopChan <- syscall.SIGINT
	time.Sleep(2 * time.Second)

	if CurrentServer != nil {
		t.Error("Server was not reset")
	}

	if !closeChan {
		t.Error("Server was not closed")
	}
}
