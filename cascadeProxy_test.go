package main

import (
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
	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving end proxy server at localhost:8082")
		err := http.ListenAndServe("localhost:8082", endProxy)
		if err != nil {
			t.Error("Error while serving end server ", err)
		}
	}()

	middleProxy := CASCADE.Run(true, "http://localhost:8082", username, password)

	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving middle proxy server at localhost:8081")
		err := http.ListenAndServe("localhost:8081", middleProxy)
		if err != nil {
			t.Error("Error while serving middle server ", err)
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
}
