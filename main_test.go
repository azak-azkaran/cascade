package main

import (
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	fmt.Println("Running: TestRun")
	utils.Init()

	config := Yaml{}
	config.HealthTime = 1
	config.Username = "foo"
	config.Password = "bar"
	config.ProxyURL = "localhost:7082"
	config.LocalPort = "7081"
	config.Log = "info"
	config.CheckAddress = "https://google.de"
	config.ConfigFile = ""
	config.VaultAddr = ""

	go Run(&config)

	time.Sleep(1 * time.Second)
	assert.NotNil(t, CurrentServer, "No Server was created")

	resp, err := utils.GetResponse("http://localhost:7081", "https://www.google.de")
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
	//args := os.Args
	//os.Args = append(os.Args, "-config=./test/config.yml")

	go main()

	time.Sleep(2 * time.Second)
	utils.Sugar.Info("calling HTTP")
	resp, err := utils.GetResponse("http://localhost:8888", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	utils.Sugar.Info("calling HTTPs")
	resp, err = utils.GetResponse("http://localhost:8888", "https://www.google.de")

	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	utils.Sugar.Info("Closing")
	stopChan <- syscall.SIGINT
	time.Sleep(2 * time.Second)

	assert.Nil(t, CurrentServer)
	assert.True(t, closeChan)
}
