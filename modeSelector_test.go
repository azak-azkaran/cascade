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
	assert.False(t, Config.OnlineCheck)

	Config.verbose = true
	Config.ProxyURL = "something"

	Config.CascadeMode = true
	DirectOverrideChan = false
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to DirectMode")
	ChangeMode(true, Config.OnlineCheck)
	assert.False(t, Config.CascadeMode)
	assert.True(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.CascadeMode = false
	DirectOverrideChan = true
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to CascadeMode")

	ChangeMode(false, Config.OnlineCheck)
	assert.True(t, Config.CascadeMode)
	assert.False(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.CascadeMode = true
	DirectOverrideChan = false
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to DirectMode")

	ChangeMode(false, true)
	assert.False(t, Config.CascadeMode)
	assert.True(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.CascadeMode = false
	DirectOverrideChan = true
	fmt.Println("Test switch from\nCascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan, " to CascadeMode")

	ChangeMode(true, true)
	assert.True(t, Config.CascadeMode)
	assert.False(t, DirectOverrideChan)
	fmt.Println("Result CascadeMode: ", Config.CascadeMode, " DirectOverrideChan: ", DirectOverrideChan)

	Config.ProxyURL = ""
}

func TestModeSelection(t *testing.T) {
	fmt.Println("Running: TestModeSelection")
	utils.Init()

	Config.verbose = true
	Config.CascadeMode = true
	Config.ProxyURL = "something"
	Config.proxyRedirectList = strings.Split("golang.org,youtube.com", ",")

	ModeSelection("https://www.asda12313.de")
	time.Sleep(1 * time.Millisecond)
	assert.False(t, DirectOverrideChan)

	ModeSelection("https://www.google.de")
	time.Sleep(1 * time.Millisecond)
	assert.True(t, DirectOverrideChan)

	Config = Yaml{}
}

func TestCreateConfig(t *testing.T) {
	fmt.Println("Running: TestCreateConfig")
	utils.Init()
	Config = Yaml{LocalPort: "8888", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "google,eclipse", Log: "info"}
	CreateConfig()

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(Config)

	assert.NotNil(t, CurrentServer)
	assert.Equal(t, len(Config.proxyRedirectList), 2)

	Config = Yaml{}
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

	Config.verbose = true
	Config.CascadeMode = true
	Config.ProxyURL = "something"
	DirectOverrideChan = false
	Config.proxyRedirectList = strings.Split("golang.org,youtube.com", ",")
	Config.DisableAutoChangeMode = true

	ModeSelection("https://www.asda12313.de")
	time.Sleep(1 * time.Millisecond)
	assert.False(t, DirectOverrideChan)

	Config = Yaml{}
}
