package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"os"
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
}

func TestModeSelection(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	cascade = false
	direct = false

	CONFIG.Verbose = true
	CONFIG.CascadeMode = true
	CONFIG.CascadeFunction = toggleCascade
	CONFIG.DirectFunction = toggleDirect

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
}
