package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
)

func TestChangeMode(t *testing.T) {
	fmt.Println("Running: TestChangeMode")
	utils.Init()
	Config := GetConfig()
	assert.False(t, Config.OnlineCheck)

	Config.verbose = true
	Config.ProxyURL = "something"

	Config.CascadeMode = true
	DirectOverrideChan = false
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to DirectMode")
	ChangeMode(true, &Config)
	assert.False(t, Config.CascadeMode)
	assert.True(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.CascadeMode = false
	DirectOverrideChan = true
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to CascadeMode")

	ChangeMode(false, &Config)
	assert.True(t, Config.CascadeMode)
	assert.False(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.CascadeMode = true
	DirectOverrideChan = false
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to DirectMode")

	Config.OnlineCheck = true
	ChangeMode(false, &Config)
	assert.False(t, Config.CascadeMode)
	assert.True(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.CascadeMode = false
	DirectOverrideChan = true
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to CascadeMode")

	Config.OnlineCheck = true
	ChangeMode(true, &Config)
	assert.True(t, Config.CascadeMode)
	assert.False(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.ProxyURL = ""
}

func TestModeSelection(t *testing.T) {
	fmt.Println("Running: TestModeSelection")
	utils.Init()

	Config := GetConfig()
	Config.verbose = true
	Config.CascadeMode = true
	Config.ProxyURL = "something"
	Config.proxyRedirectList = strings.Split("golang.org,youtube.com", ",")
	Config.CheckAddress = "https://www.asda12313.de"
	Config.OnlineCheck = false
	ModeSelection(&Config)
	time.Sleep(1 * time.Millisecond)
	assert.False(t, DirectOverrideChan)

	Config.CheckAddress = "https://www.google.de"
	ModeSelection(&Config)
	time.Sleep(1 * time.Millisecond)
	assert.True(t, DirectOverrideChan)

	Config = Yaml{}
}

func TestCreateConfig(t *testing.T) {
	fmt.Println("Running: TestCreateConfig")
	utils.Init()
	Config := Yaml{LocalPort: "8888", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "google,eclipse", Log: "info"}
	conf := CreateConfig(&Config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(conf)

	assert.NotNil(t, CurrentServer)
	assert.Equal(t, len(Config.proxyRedirectList), 2)
}

func TestHandleCustomProxies(t *testing.T) {
	fmt.Println("Running: TestHandleCustomProxies")
	utils.Init()
	list := strings.Split("eclipse,google->test:8888,azure->", ",")
	HandleCustomProxies(list)

	val, in := HostList.Get("")
	assert.True(t, in)

	value := val.(hostConfig)
	assert.True(t, value.reg.MatchString("eclipse2017.nasa.gov"))
	assert.True(t, in)
	assert.False(t, value.proxyAddr != "")

	val, in = HostList.Get("test:8888")
	assert.True(t, in)

	value = val.(hostConfig)
	assert.True(t, value.reg.MatchString("www.google.de"))
	assert.Equal(t, strings.Compare(value.proxyAddr, "http://test:8888"), 0)

	val, in = HostList.Get("")
	assert.True(t, in)

	value = val.(hostConfig)

	assert.True(t, value.reg.MatchString("https://azure.microsoft.com/en-us/"))
	assert.False(t, value.proxyAddr != "")
}

func TestDisableAutoChangeMode(t *testing.T) {
	fmt.Println("Running: TestDisableAutoChangeMode")
	utils.Init()

	Config := GetConfig()
	Config.verbose = true
	Config.CascadeMode = true
	Config.ProxyURL = "something"
	DirectOverrideChan = false
	Config.proxyRedirectList = strings.Split("golang.org,youtube.com", ",")
	Config.DisableAutoChangeMode = true
	Config.CheckAddress = "https://www.asda12313.de"

	ModeSelection(&Config)
	time.Sleep(1 * time.Millisecond)
	assert.False(t, DirectOverrideChan)
}
