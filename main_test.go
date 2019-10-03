package main

import (
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGetConf(t *testing.T) {
	fmt.Println("Running: TestGetConf")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	conf, _ := GetConf("./test/test.yml")

	assert.Equal(t, "TestHealth", conf.CheckAddress)
	assert.Equal(t, "8888", conf.LocalPort)
	assert.Equal(t, "TestHost", conf.ProxyURL)
	assert.Equal(t, "TestPassword", conf.Password)
	assert.Equal(t, "TestUser", conf.Username)
	assert.Equal(t, int64(5), conf.HealthTime)

	conf, err := GetConf("noname.yaml")
	assert.Error(t, err)
	assert.Nil(t, conf, "Error could read YAML but should not be able to be")
}

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
