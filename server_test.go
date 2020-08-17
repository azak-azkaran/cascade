package main

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateServer(t *testing.T) {
	fmt.Println("Running: TestCreateServer")
	utils.Init()
	test_config := Yaml{LocalPort: "7082", Verbose: true}
	conf := CreateConfig(&test_config)
	testServer := CreateServer(conf) //CASCADE.Run(true, "", "", ""), "localhost", "7082")
	go func() {
		err := testServer.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error("Error while running server", err)
		}

	}()

	time.Sleep(1 * time.Second)
	assert.NotNil(t, CurrentServer)
	resp, err := utils.GetResponse("http://localhost:7082", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = testServer.Shutdown(context.TODO())
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	resp, err = utils.GetResponse("http://localhost:7082", "http://www.google.de")
	assert.Error(t, err)

	//err = utils.Sugar.Sync()
	//assert.NoError(t, err)
}

func TestRunServer(t *testing.T) {
	fmt.Println("Running: TestRunServer")
	utils.Init()
	test_config := Yaml{LocalPort: "7082", Verbose: true}
	conf := CreateConfig(&test_config)
	testServer := CreateServer(conf)
	require.False(t, running)
	require.NotNil(t, CurrentServer)

	RunServer()
	time.Sleep(1 * time.Second)
	assert.True(t, running)

	resp, err := utils.GetResponse("http://localhost:7082", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = testServer.Shutdown(context.TODO())
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	CurrentServer = nil
}

func TestShutdownCurrentServer(t *testing.T) {
	fmt.Println("Running: TestShutdownCurrentServer")
	utils.Init()
	test_config := Yaml{LocalPort: "7082", Verbose: true}
	conf := CreateConfig(&test_config)
	CreateServer(conf)
	assert.False(t, running)

	ShutdownCurrentServer()
	time.Sleep(1 * time.Second)
	assert.False(t, running)
	assert.Nil(t, CurrentServer)
}

func TestCreateBrokenServer(t *testing.T) {
	fmt.Println("Running: TestCreateBrokenServer")
	utils.Init()
	Config := Yaml{LocalPort: "7082",
		CheckAddress: "https://www.google.de",
		HealthTime:   5,
		HostList:     "golang.org,youtube.com",
		Log:          "info",
		CascadeMode:  true,
	}
	conf := CreateConfig(&Config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(conf)

	RunServer()
	time.Sleep(1 * time.Second)
	assert.True(t, running)
	assert.Len(t, Config.proxyRedirectList, 2)

	_, err := utils.GetResponse("http://localhost:7082", "https://www.google.de")
	assert.Error(t, err)

	resp, err := utils.GetResponse("http://localhost:7082", "http://golang.org/doc/")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ShutdownCurrentServer()
	time.Sleep(1 * time.Second)
	assert.False(t, running)
	assert.Nil(t, CurrentServer)
}

func TestRestRequest(t *testing.T) {
	fmt.Println("Running: TestRestRequest")
	utils.Init()
	test_config := Yaml{LocalPort: "7082", Verbose: true}
	conf := CreateConfig(&test_config)
	testServer := CreateServer(conf)

	RunServer()
	time.Sleep(1 * time.Second)

	assert.True(t, running, "Server was not started")

	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)

	fmt.Println("Direct Config Call")
	resp, err := client.Get("http://localhost:7082/config")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	fmt.Println("Proxied Config Call")
	resp, err = utils.GetResponse("http://localhost:7082", "http://localhost:7082/config")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	fmt.Println("Proxied Google Call")
	resp, err = utils.GetResponse("http://localhost:7082", "https://www.github.com")
	assert.NoError(t, err, "Error while client request over proxy server", err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Error while client request over proxy server", resp)

	err = testServer.Shutdown(context.TODO())
	assert.NoError(t, err, "Error while shutdown")

	err = CurrentServer.Shutdown(context.TODO())
	assert.NoError(t, err, "Error while shutdown")
}

func TestRestServerLateCreation(t *testing.T) {
	fmt.Println("Running: TestRestServerNegative")
	utils.Init()
	CurrentServer = nil
	running = false

	RunServer()
	time.Sleep(3 * time.Second)
	assert.False(t, running)
	test_config := Yaml{LocalPort: "7082", Verbose: true}
	conf := CreateConfig(&test_config)
	testServer := CreateServer(conf)
	assert.Eventually(t, func() bool { return running }, 5*time.Second, 10*time.Millisecond)
	assert.True(t, running)
	assert.NotNil(t, CurrentServer)

	err := testServer.Shutdown(context.TODO())
	assert.NoError(t, err, "Error while shutdown")
}
