package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRestRouter_GetConfigWithProxy(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigWithProxy")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	endProxy := DIRECT.Run(true)

	n := negroni.Classic()

	CreateRestEndpoint("localhost", "8081", false)
	n.Use(negroni.Wrap(RestRouter))
	n.UseHandler(endProxy)
	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: n,
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)
	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestRestRouter_GetConfigWithNegroni(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigWithNegroni")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	CreateConfig("8082", "", "", "", "https://www.google.de", 5, "golang.org,youtube.com", "info")
	CreateRestEndpoint("localhost", "8081", true)
	n := negroni.Classic()
	n.Use(negroni.Wrap(RestRouter))

	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: n,
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)

	client, err := utils.GetClient("http://localhost:8081", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:8081/config")
	assert.NoError(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)

	resp, err = client.Get("http://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestRestRouter_GetConfigOnlyMux(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigOnlyMux")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CreateConfig("8082", "", "", "", "https://www.google.de", 5, "golang.org,youtube.com", "info")
	CreateRestEndpoint("localhost", "8081", true)

	//n.UseFunc(HandleConfig)

	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: RestRouter,
	}

	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)

	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:8081/config")
	assert.NoError(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)

	decoder := json.NewDecoder(resp.Body)
	var decodedConfig Yaml
	assert.NoError(t, decoder.Decode(&decodedConfig))

	assert.Equal(t, Config.CascadeMode, decodedConfig.CascadeMode)
	assert.Equal(t, Config.HealthTime, decodedConfig.HealthTime)
	assert.Equal(t, Config.CheckAddress, decodedConfig.CheckAddress)
	assert.Equal(t, Config.HostList, decodedConfig.HostList)
	assert.Equal(t, Config.LocalPort, decodedConfig.LocalPort)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
}
