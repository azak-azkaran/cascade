package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.com")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}
}

func TestRestRouter_AddRedirect(t *testing.T) {
	fmt.Println("Running: TestRestRouter_AddReddirect")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	Config = Yaml{LocalPort: "8082", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "golang.org,youtube.com", Log: "info", ConfigFile: "test.yml"}
	CreateConfig()

	r := gin.Default()
	endServer := &http.Server{
		Addr:    "localhost:8081",
		Handler: r,
	}

	r.POST("/add", addRedirectFunc)
	go endServer.ListenAndServe()

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

	resp, err := client.Post("http://localhost:8081/add", "application/json", &buf)
	assert.NoError(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	assert.NotNil(t, resp.Body)

	decoder := json.NewDecoder(resp.Body)
	var responesMessage AddRedirect
	assert.NoError(t, decoder.Decode(&responesMessage))
	assert.Equal(t, jsonRequest.Proxy, responesMessage.Proxy)
	assert.Equal(t, jsonRequest.Address, responesMessage.Address)

	utils.Info.Println("Config: ", Config)

	value, available := HostList.Get(jsonRequest.Proxy)
	assert.True(t, available)
	assert.NotNil(t, value)

	config, err := GetConf("test.yml")
	assert.NoError(t, err)
	assert.Equal(t, config.LocalPort, Config.LocalPort)
	assert.Equal(t, config.HealthTime, Config.HealthTime)
	assert.Equal(t, config.HostList, Config.HostList)
	assert.Equal(t, config.CheckAddress, Config.CheckAddress)
	assert.NotEqual(t, config.HostList, "golang.org,youtube.com")

	buf = *bytes.NewBufferString("hallo")
	resp, err = client.Post("http://localhost:8081/add", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	err = os.Remove("test.yml")
	assert.NoError(t, err)
	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}
}

func TestRestRouter_ChangeOnlineCheck(t *testing.T) {
	fmt.Println("Running: TestRestRouter_ChangeOnlineCheck")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	Config = Yaml{OnlineCheck: false}
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

	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:8081/getOnlineCheck")
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
	resp, err = client.Post("http://localhost:8081/setOnlineCheck", "application/json", &buf)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(bodyBytes, &jsonRequest))
	assert.True(t, jsonRequest.OnlineCheck)

	buf = *bytes.NewBufferString("[ hallo ]")
	resp, err = client.Post("http://localhost:8081/setOnlineCheck", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	resp, err = client.Get("http://localhost:8081/getOnlineCheck")
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.NotNil(t, resp.Body)

	decoder = json.NewDecoder(resp.Body)
	assert.NoError(t, decoder.Decode(&decodedBool))
	assert.True(t, decodedBool)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}
}

func TestRestRouter_DisableAutomaticChange(t *testing.T) {
	fmt.Println("Running: TestRestRouter_DisableAutomaticChange")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	Config = Yaml{DisableAutoChangeMode: false, ProxyURL: "http://localhost", Log: "DEBUG", OnlineCheck: false}
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
	client, err := utils.GetClient("", 2)
	assert.NoError(t, err)

	resp, err := client.Get("http://localhost:8081/getAutoMode")
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
	resp, err = client.Post("http://localhost:8081/setAutoMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)
	defer resp.Body.Close()

	decoder = json.NewDecoder(resp.Body)
	assert.NoError(t, decoder.Decode(&jsonRequest))
	assert.False(t, jsonRequest.AutoChangeMode)
	assert.True(t, Config.DisableAutoChangeMode)
	assert.True(t, Config.CascadeMode)
	assert.False(t, Config.OnlineCheck)

	ModeSelection("https://www.asda12313.de")
	time.Sleep(1 * time.Millisecond)
	assert.True(t, Config.CascadeMode)
	assert.True(t, Config.DisableAutoChangeMode)

	buf = *bytes.NewBufferString("hallo")
	resp, err = client.Post("http://localhost:8081/setAutoMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}

}
func TestRestRouter_ChangeCascadeMode(t *testing.T) {
	fmt.Println("Running: TestRestRouter_ChangeCascadeMode")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	Config = Yaml{
		DisableAutoChangeMode: true,
		ProxyURL:              "http://localhost",
		Log:                   "DEBUG",
		OnlineCheck:           false,
		CascadeMode:           false,
	}

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
	assert.False(t, DirectOverrideChan)
	assert.False(t, Config.CascadeMode)

	cascadeModeReq := SetCascadeModeRequest{
		CascadeMode: true,
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err = encoder.Encode(&cascadeModeReq)
	assert.NoError(t, err)
	resp, err := client.Post("http://localhost:8081/setCascadeMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, resp.Body)
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	assert.NoError(t, decoder.Decode(&cascadeModeReq))
	assert.True(t, cascadeModeReq.CascadeMode)
	assert.True(t, Config.DisableAutoChangeMode)
	assert.True(t, Config.CascadeMode)
	assert.False(t, DirectOverrideChan)

	cascadeModeReq.CascadeMode = false
	err = encoder.Encode(&cascadeModeReq)
	assert.NoError(t, err)
	resp, err = client.Post("http://localhost:8081/setCascadeMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, cascadeModeReq.CascadeMode)
	assert.False(t, Config.CascadeMode)
	assert.True(t, DirectOverrideChan)

	buf = *bytes.NewBufferString("hallo")
	resp, err = client.Post("http://localhost:8081/setCascadeMode", "application/json", &buf)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	err = endServer.Shutdown(context.Background())
	assert.NoError(t, err)
	Config = Yaml{}
}
