package main

import (
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"strings"
	"testing"
	"time"
)

func TestChangeMode(t *testing.T) {
	fmt.Println("Running: TestChangeMode")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	Config.CascadeMode = true
	Config.ProxyURL = "something"
	ChangeMode(true)
	if Config.CascadeMode {
		t.Error("Mode was not changed")
	}

	time.Sleep(1 * time.Second)
	if !DirectOverrideChan {
		t.Error("DirectOverride is active")
	}

	ChangeMode(false)
	if !Config.CascadeMode {
		t.Error("Mode was not changed")
	}
	time.Sleep(1 * time.Second)
	if DirectOverrideChan {
		t.Error("DirectOverride is not active")
	}
	Config.ProxyURL = ""
}

func TestModeSelection(t *testing.T) {
	fmt.Println("Running: TestModeSelection")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	Config.verbose = true
	Config.CascadeMode = true
	Config.ProxyURL = "something"
	Config.proxyRedirectList = strings.Split("golang.org,youtube.com", ",")

	ModeSelection("https://www.asda12313.de")
	time.Sleep(1 * time.Second)
	if DirectOverrideChan {
		t.Error("DirectOverride is active")
	}

	ModeSelection("https://www.google.de")
	time.Sleep(1 * time.Second)
	if !DirectOverrideChan {
		t.Error("DirectOverride is not active")
	}

	Config = Yaml{}
}

func TestCreateConfig(t *testing.T) {
	fmt.Println("Running: TestCreateConfig")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	Config = Yaml{}
	CreateConfig("8888", "", "", "", "https://www.google.de", 5, "google,eclipse")

	if CurrentServer == nil {
		t.Error("Server was not created")
	}

	if len(Config.proxyRedirectList) != 2 {
		t.Error("SkipHosts was not split correctly")
	}

	Config = Yaml{}
}

func TestHandleCustomProxies(t *testing.T) {
	fmt.Println("Running: TestHandleCustomProxies")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	list := strings.Split("eclipse,google->test:8888,azure->", ",")
	HandleCustomProxies(list)

	val, in := HostList.Get("")
	if !in {
		t.Error("Proxy redirect to eclipse not added")
	}

	value := val.(HostConfig)
	if in && !value.reg.MatchString("eclipse2017.nasa.gov") {
		t.Error("Proxy redirect regex does not match: ", value.regString)
	}
	if value.proxyAddr != "" && in {
		t.Error("Proxy redirect regex does not match")

	}

	val, in = HostList.Get("test:8888")
	if !in {
		t.Error("Proxy redirect to google not added: ", value.regString)
	}

	if in {
		value = val.(HostConfig)
		if !value.reg.MatchString("www.google.de") {
			t.Error("Proxy redirect regex does not match: ", value.regString)
		}
		if strings.Compare(value.proxyAddr, "http://test:8888") != 0 {
			t.Error("Proxy redirect address does not match")
		}
	}

	val, in = HostList.Get("")
	if !in {
		t.Error("Proxy redirect to azure not added: ", value.regString)
	}

	if in {
		value = val.(HostConfig)

		if !value.reg.MatchString("https://azure.microsoft.com/en-us/") {
			t.Error("Proxy redirect regex does not match: ", value.regString)
		}
		if value.proxyAddr != "" {
			t.Error("Proxy redirect address does not match")
		}
	}
}
