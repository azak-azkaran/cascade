package main

import (
	"context"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/elazarl/goproxy"
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
	endProxy := goproxy.NewProxyHttpServer()
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

	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
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
