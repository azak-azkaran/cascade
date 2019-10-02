package main

import (
	"context"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy/ext/auth"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestCascadeProxy_Run(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	username, password := "foo", "bar"

	// start end proxy server
	endProxy := DIRECT.Run(true)
	endProxy.Verbose = true
	auth.ProxyBasic(endProxy, "my_realm", func(user, pwd string) bool {
		return user == username && password == pwd
	})
	var endServer *http.Server
	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving end proxy server at localhost:8082")
		endServer = &http.Server{
			Addr:    "localhost:8082",
			Handler: endProxy,
		}
		err := endServer.ListenAndServe()
		if err == nil {
			t.Error("Error shutdown should always return error", err)
		}
	}()

	middleProxy := CASCADE.Run(true, "http://localhost:8082", username, password)
	var middleServer *http.Server

	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		if err == nil {
			t.Error("Error shutdown should always return error", err)
		}
	}()

	utils.Info.Println("waiting for running")
	time.Sleep(1 * time.Second)

	utils.Info.Println("Start http Test")
	client, err := utils.GetClient("http://localhost:8081", 2)
	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	request, _ := http.NewRequest("GET", "http://google.de/", nil)
	resp, err := client.Do(request)

	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp)
	}

	utils.Info.Println("Start https Test")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp)
	}

	_, err = utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err == nil {
		t.Error("Error direct connection did also work", err)
	}

	err = middleServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down middle Server, ", err)
	}

	err = endServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down end Server, ", err)
	}
	time.Sleep(1 * time.Second)
}

func TestAddDirectConnection(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	middleProxy := CASCADE.Run(true, "http://localhost:8082", "", "")
	var middleServer *http.Server

	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		if err == nil {
			t.Error("Error shutdown should always return error", err)
		}
	}()

	time.Sleep(1 * time.Second)
	_, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err == nil {
		t.Error("Error while requesting google", err)
	}

	Config.LocalPort = "8801"
	AddDirectConnection("google")
	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	ClearHostList()
	_, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err == nil {
		t.Error("Error while requesting google", err)
	}

	err = middleServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down middle Server, ", err)
	}
}

func TestAddDifferentProxyConnection(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	username, password := "foo", "bar"
	endProxy := DIRECT.Run(true)
	endProxy.Verbose = true
	var endServer *http.Server
	auth.ProxyBasic(endProxy, "my_realm", func(user, pwd string) bool {
		return user == username && password == pwd
	})
	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving end proxy server at localhost:8082")
		endServer = &http.Server{
			Addr:    "localhost:8082",
			Handler: endProxy,
		}
		err := endServer.ListenAndServe()
		if err == nil {
			t.Error("Error shutdown should always return error", err)
		}
	}()

	middleProxy := CASCADE.Run(true, "http://localhost:8083", username, password)
	var middleServer *http.Server

	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		if err == nil {
			t.Error("Error shutdown should always return error", err)
		}
	}()

	time.Sleep(1 * time.Second)
	Config.LocalPort = "8801"
	utils.Info.Println("starting HTTPS test to fail")
	_, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err == nil {
		t.Error("Error while requesting google", err)
	}
	utils.Info.Println("starting HTTP test to fail")
	resp, err := utils.GetResponse("http://localhost:8081", "http://www.google.de")
	if resp == nil || resp.StatusCode != 500 {
		t.Error("Error while requesting google", err)
	}

	AddDifferentProxyConnection("google", "http://localhost:8082")
	utils.Info.Println("starting HTTPS test to work")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	utils.Info.Println("starting HTTP test")
	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	AddDifferentProxyConnection("google", "")
	utils.Info.Println("starting HTTPS test to work")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	utils.Info.Println("starting HTTP test")
	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
		if err != nil {
			t.Error("Error while requesting google", err)
		}
		if resp == nil || resp.StatusCode != 200 {
			t.Error("Google was not available")
		}

		utils.Info.Println("starting HTTP test")
		resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
		if err != nil {
			t.Error("Error while requesting google", err)
		}
		if resp == nil || resp.StatusCode != 200 {
			t.Error("Google was not available")
		}
		t.Error("Google was not available")
	}

	AddDifferentProxyConnection("steam", "")
	if !HostList.Has("") {
		t.Error("direct not in HostList: ", HostList.Keys())
	}

	ClearHostList()
	_, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err == nil {
		t.Error("Error while requesting google", err)
	}

	err = middleServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down middle Server, ", err)
	}
	err = endServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down middle Server, ", err)
	}
}

func TestCascadeProxy_ModeSwitch(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	middleProxy := CASCADE.Run(true, "http://localhost:8083", "", "")
	var middleServer *http.Server

	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		if err == nil {
			t.Error("Error shutdown should always return error", err)
		}
	}()

	time.Sleep(1 * time.Millisecond)
	_, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err == nil {
		t.Error("Error while requesting google", err)
	}

	utils.Info.Println("writing to DirectOverrideChan")
	DirectOverrideChan = true
	utils.Info.Println("Testing direct Override")
	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting google in directOverride", err)
	}

	if resp == nil || resp.StatusCode != 200 {
		t.Error("Error while requesting google in directOverride", resp)
	}

	err = middleServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down middle Server, ", err)
	}
	DirectOverrideChan = false
}
