package main

import (
	"context"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/elazarl/goproxy/ext/auth"
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
	client , err:= utils.GetClient("http://localhost:8081")
	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	request, err := http.NewRequest("GET", "http://google.de/", nil)
	resp, err := client.Do( request)

	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp.Status)
	}

	utils.Info.Println("Start https Test")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp.Status)
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

func TestHandleDirect(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	resp, err := utils.GetResponse("", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	req, err := http.NewRequest("GET", "https://www.google.de", nil)
	if err != nil {
		t.Error("Error while creating request to google", err)
	}

	req, resp = HandleDirect(req, nil)
	if resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

	req, err = http.NewRequest("GET", "http://www.google.de", nil)
	req, resp = HandleDirect(req, nil)
	if resp.StatusCode != 200 {
		t.Error("Google was not available")
	}
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
	if err == nil{
		t.Error("Error while requesting google", err)
	}

	AddDirectConnection(middleProxy, "google")
	resp, err := utils.GetResponse("http://localhost:8081", "http://www.google.de")
	if err != nil {
		t.Error("Error while requesting google", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Google was not available")
	}

err = middleServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down middle Server, ", err)
	}
}
