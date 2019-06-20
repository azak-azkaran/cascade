package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"strings"
	"testing"
	"time"
)

var cascade bool
var direct bool

func toggleCascade() {
	cascade = !cascade
}

func toggleDirect() {
	direct = !direct
}

func TestChangeMode(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	cascade = false
	direct = false

	CONFIG.CascadeMode = true
	CONFIG.CascadeFunction = toggleCascade
	CONFIG.DirectFunction = toggleDirect
	CONFIG.ProxyURL = "something"
	ChangeMode(true)
	if CONFIG.CascadeMode {
		t.Error("Mode was not changed")
	}

	time.Sleep(1 * time.Second)
	if !direct {
		t.Error("direct function was not called")
	}

	ChangeMode(false)
	if !CONFIG.CascadeMode {
		t.Error("Mode was not changed")
	}
	time.Sleep(1 * time.Second)
	if !cascade {
		t.Error("cascade function was not called")
	}
	CONFIG.CascadeFunction = nil
	CONFIG.DirectFunction = nil
	CONFIG.ProxyURL = ""
}

func TestModeSelection(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	cascade = false
	direct = false

	CONFIG.Verbose = true
	CONFIG.CascadeMode = true
	CONFIG.CascadeFunction = toggleCascade
	CONFIG.DirectFunction = toggleDirect
	CONFIG.ProxyURL = "something"
	CONFIG.SkipCascadeHosts = strings.Split("golang.org,youtube.com", ",")

	ModeSelection("https://www.asda12313.de")
	time.Sleep(1 * time.Second)
	if !cascade {
		t.Error("cascade function was not called")
	}

	ModeSelection("https://www.google.de")
	time.Sleep(1 * time.Second)
	if !direct {
		t.Error("direct function was not called")
	}

	CONFIG = config{}
}

func TestCreateConfig(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CONFIG = config{}
	CreateConfig("8888", "", "", "", "https://www.google.de", 5, "google,eclipse")

	if CONFIG.CascadeFunction == nil {
		t.Error("Cascade function was not created")
	}

	if len(CONFIG.SkipCascadeHosts) != 2 {
		t.Error("SkipHosts was not split correctly")
	}

	CONFIG = config{}
}

func TestHandleCustomProxies(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	list := strings.Split("eclipse,google->test:8888,azure->", ",")
	HandleCustomProxies(list)

	val, in := HostList.Get("eclipse")
	if !in {
		t.Error("Proxy redirect to eclipse not added")
	}

	value := val.(HostConfig)
	if in && !value.reg.MatchString("eclipse2017.nasa.gov") {
		t.Error("Proxy redirect regex does not match")
	}
	if value.proxyAddr != "" && in {
		t.Error("Proxy redirect regex does not match")

	}

	val, in = HostList.Get("google")
	if !in {
		t.Error("Proxy redirect to google not added")
	}

	if in {
		value = val.(HostConfig)

		if !value.reg.MatchString("www.google.de") {
			t.Error("Proxy redirect regex does not match")
		}
		if strings.Compare(value.proxyAddr, "test:8888") != 0 {
			t.Error("Proxy redirect address does not match")
		}
	}

	val, in = HostList.Get("azure")
	if !in {
		t.Error("Proxy redirect to azure not added")
	}

	if in {
		value = val.(HostConfig)

		if !value.reg.MatchString("https://azure.microsoft.com/en-us/") {
			t.Error("Proxy redirect regex does not match")
		}
		if value.proxyAddr != "" {
			t.Error("Proxy redirect address does not match")
		}
	}

}
