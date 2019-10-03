package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRestRouter_RouteToOtherLocalhost(t *testing.T) {
	fmt.Println("Running: TestRestRouter_RouteToOtherLocalhost")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	endProxy := DIRECT.Run(true)

	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: ConfigureRouter(endProxy, "localhost", true),
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	r := gin.Default()
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})
	otherServer := &http.Server{
		Addr:    "localhost:3000",
		Handler: r,
	}
	go func() {
		err := otherServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)

	resp, err := utils.GetResponse("", "http://localhost:3000/someJSON")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = utils.GetResponse("http://localhost:8081", "http://localhost:3000/someJSON")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)

	err = otherServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestRestRouter_GetConfigWithProxy(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigWithProxy")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	endProxy := DIRECT.Run(true)

	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: ConfigureRouter(endProxy, "localhost", true),
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

func TestRestRouter_GetConfigWithMiddleware(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigWithMiddleware")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	Config = Yaml{LocalPort: "8082", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "golang.org,youtube.com", Log: "info"}
	CreateConfig()

	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: ConfigureRouter(DIRECT.Run(true), "localhost", true),
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
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestRestRouter_GetConfigOnlyMux(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigOnlyMux")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	Config = Yaml{LocalPort: "8082", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "golang.org,youtube.com", Log: "info"}
	CreateConfig()

	//n.UseFunc(HandleConfig)

	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: ConfigureRouter(DIRECT.Run(true), "localhost", true),
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
