package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestRouter_RouteToOtherLocalhost(t *testing.T) {
	fmt.Println("Running: TestRestRouter_RouteToOtherLocalhost")
	utils.Init()
	endProxy := DIRECT.Run(true)

	endServer := &http.Server{
		Addr:    "localhost:7081",
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

	resp, err = utils.GetResponse("http://localhost:7081", "http://localhost:3000/someJSON")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)

	err = otherServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestRestRouter_GetConfigWithProxy(t *testing.T) {
	fmt.Println("Running: TestRestRouter_GetConfigWithProxy")
	utils.Init()

	Config := Yaml{LocalPort: "7082", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "golang.org,youtube.com", Log: "info"}
	conf := CreateConfig(&Config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(conf)

	endServer := &http.Server{
		Addr:    "localhost:7081",
		Handler: ConfigureRouter(DIRECT.Run(true), "localhost", true),
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)

	client, err := utils.GetClient("http://localhost:7081", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:7081/config")
	assert.NoError(t, err)

	decoder := json.NewDecoder(resp.Body)
	var decodedConfig Yaml
	assert.NoError(t, decoder.Decode(&decodedConfig))

	assert.Equal(t, Config.CascadeMode, decodedConfig.CascadeMode)
	assert.Equal(t, Config.HealthTime, decodedConfig.HealthTime)
	assert.Equal(t, Config.CheckAddress, decodedConfig.CheckAddress)
	assert.Equal(t, Config.HostList, decodedConfig.HostList)
	assert.Equal(t, Config.LocalPort, decodedConfig.LocalPort)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)

	resp, err = utils.GetResponse("", "http://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = utils.GetResponse("http://localhost:7081", "http://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = utils.GetResponse("http://localhost:7081", "https://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}
}

func TestRestRouter_AddRedirect(t *testing.T) {
	fmt.Println("Running: TestRestRouter_AddReddirect")
	utils.Init()
	Config := Yaml{LocalPort: "7082", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "golang.org,youtube.com", Log: "info", ConfigFile: "test/config.yml"}
	conf := CreateConfig(&Config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(conf)

	r := gin.Default()
	endServer := &http.Server{
		Addr:    "localhost:7081",
		Handler: r,
	}

	r.POST("/add", addRedirectFunc)
	go func() {
		err := endServer.ListenAndServe()
		require.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)
	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)
	jsonRequest := AddRedirect{
		Proxy:   "test.dlh.de",
		Address: "10.20.10.20",
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err = encoder.Encode(&jsonRequest)
	assert.NoError(t, err)
	//var jsonRequest = []byte(`{"proxy":"test.dlh.de", "address":"10.20.10.20"}`)

	resp, err := client.Post("http://localhost:7081/add", "application/json", &buf)
	assert.NoError(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	require.NotNil(t, resp.Body)

	decoder := json.NewDecoder(resp.Body)
	var responesMessage AddRedirect
	assert.NoError(t, decoder.Decode(&responesMessage))
	assert.Equal(t, jsonRequest.Proxy, responesMessage.Proxy)
	assert.Equal(t, jsonRequest.Address, responesMessage.Address)

	utils.Sugar.Info("Config: ", Config)

	value, available := HostList.Get(jsonRequest.Proxy)
	assert.True(t, available)
	assert.NotNil(t, value)

	buf = *bytes.NewBufferString("hallo")
	resp, err = client.Post("http://localhost:7081/add", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)

	err = CurrentServer.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestRestRouter_ChangeOnlineCheck(t *testing.T) {
	fmt.Println("Running: TestRestRouter_ChangeOnlineCheck")
	utils.Init()

	Config := Yaml{OnlineCheck: false}
	conf := CreateConfig(&Config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(conf)

	endServer := &http.Server{
		Addr:    "localhost:7081",
		Handler: ConfigureRouter(DIRECT.Run(true), "localhost", true),
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)

	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:7081/getOnlineCheck")
	assert.NoError(t, err)
	assert.NotNil(t, resp.Body)
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var decodedBool bool
	assert.NoError(t, decoder.Decode(&decodedBool))
	assert.False(t, decodedBool)

	jsonRequest := SetOnlineCheckRequest{
		OnlineCheck: true,
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err = encoder.Encode(&jsonRequest)
	assert.NoError(t, err)
	resp, err = client.Post("http://localhost:7081/setOnlineCheck", "application/json", &buf)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(bodyBytes, &jsonRequest))
	assert.True(t, jsonRequest.OnlineCheck)

	buf = *bytes.NewBufferString("[ hallo ]")
	resp, err = client.Post("http://localhost:7081/setOnlineCheck", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = client.Get("http://localhost:7081/getOnlineCheck")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.NotNil(t, resp.Body)

	decoder = json.NewDecoder(resp.Body)
	assert.NoError(t, decoder.Decode(&decodedBool))
	assert.True(t, decodedBool)

	err = CurrentServer.Shutdown(context.Background())
	assert.NoError(t, err)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}
}

func TestRestRouter_DisableAutomaticChange(t *testing.T) {
	fmt.Println("Running: TestRestRouter_DisableAutomaticChange")
	utils.Init()

	Config := Yaml{DisableAutoChangeMode: false,
		ProxyURL:    "http://localhost",
		Log:         "DEBUG",
		OnlineCheck: false,
		CascadeMode: true}
	conf := CreateConfig(&Config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(conf)

	endServer := &http.Server{
		Addr:    "localhost:7081",
		Handler: ConfigureRouter(DIRECT.Run(true), "localhost", true),
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)
	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:7081/getAutoMode")
	assert.NoError(t, err)
	assert.NotNil(t, resp.Body)
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var decodedBool bool
	assert.NoError(t, decoder.Decode(&decodedBool))
	assert.True(t, decodedBool)

	jsonRequest := SetDisableAutoChangeModeRequest{
		AutoChangeMode: false,
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err = encoder.Encode(&jsonRequest)
	assert.NoError(t, err)
	resp, err = client.Post("http://localhost:7081/setAutoMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)
	defer resp.Body.Close()

	decoder = json.NewDecoder(resp.Body)
	assert.NoError(t, decoder.Decode(&jsonRequest))
	assert.False(t, jsonRequest.AutoChangeMode)

	conf = GetConfig()
	assert.True(t, conf.DisableAutoChangeMode)
	assert.True(t, conf.CascadeMode)
	assert.False(t, conf.OnlineCheck)

	Config.CheckAddress = "https://www.asda12313.de"

	ModeSelection(&Config)
	time.Sleep(1 * time.Millisecond)

	conf = GetConfig()
	assert.True(t, conf.CascadeMode)
	assert.True(t, conf.DisableAutoChangeMode)

	buf = *bytes.NewBufferString("hallo")
	resp, err = client.Post("http://localhost:7081/setAutoMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	err = CurrentServer.Shutdown(context.Background())
	assert.NoError(t, err)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}

}
func TestRestRouter_ChangeCascadeMode(t *testing.T) {
	fmt.Println("Running: TestRestRouter_ChangeCascadeMode")
	utils.Init()

	Config := Yaml{
		DisableAutoChangeMode: true,
		ProxyURL:              "http://localhost",
		Log:                   "DEBUG",
		OnlineCheck:           false,
		CascadeMode:           false,
	}
	conf := CreateConfig(&Config)

	endServer := &http.Server{
		Addr:    "localhost:7081",
		Handler: ConfigureRouter(DIRECT.Run(true), "localhost", true),
	}
	go func() {
		err := endServer.ListenAndServe()
		assert.Equal(t, http.ErrServerClosed, err)
	}()

	time.Sleep(1 * time.Second)
	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)
	assert.False(t, Config.CascadeMode)

	cascadeModeReq := SetCascadeModeRequest{
		CascadeMode: true,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err = encoder.Encode(&cascadeModeReq)
	assert.NoError(t, err)
	resp, err := client.Post("http://localhost:7081/setCascadeMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	assert.NoError(t, decoder.Decode(&cascadeModeReq))
	assert.True(t, cascadeModeReq.CascadeMode)

	conf = GetConfig()
	assert.True(t, conf.DisableAutoChangeMode)
	assert.True(t, conf.CascadeMode)

	cascadeModeReq.CascadeMode = false
	err = encoder.Encode(&cascadeModeReq)
	assert.NoError(t, err)
	resp, err = client.Post("http://localhost:7081/setCascadeMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, cascadeModeReq.CascadeMode)

	conf = GetConfig()
	assert.False(t, conf.CascadeMode)

	buf = *bytes.NewBufferString("hallo")
	resp, err = client.Post("http://localhost:7081/setCascadeMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
}
