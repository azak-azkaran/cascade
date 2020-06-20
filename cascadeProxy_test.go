package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy/ext/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCascadeProxy_Run(t *testing.T) {
	utils.Init()
	username, password := "foo", "bar"

	// start end proxy server
	endProxy := DIRECT.Run(true)
	endProxy.Verbose = true
	auth.ProxyBasic(endProxy, "my_realm", func(user, pwd string) bool {
		return user == username && password == pwd
	})
	var endServer *http.Server
	go func() {
		utils.Sugar.Info("serving end proxy server at localhost:8082")
		endServer = &http.Server{
			Addr:    "localhost:8082",
			Handler: endProxy,
		}
		err := endServer.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	middleProxy := CASCADE.Run(true, "http://localhost:8082", username, password)
	var middleServer *http.Server

	go func() {
		utils.Sugar.Info("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	utils.Sugar.Info("waiting for running")
	time.Sleep(1 * time.Second)

	utils.Sugar.Info("Start http Test")
	client, err := utils.GetClient("http://localhost:8081", 2)
	assert.NoError(t, err)

	request, _ := http.NewRequest("GET", "http://google.de/", nil)
	resp, err := client.Do(request)
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	utils.Sugar.Info("Start https Test")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	_, err = utils.GetResponse("http://localhost:8082", "https://www.google.de")
	assert.NoError(t, err)

	err = middleServer.Shutdown(context.TODO())
	assert.NoError(t, err)

	err = endServer.Shutdown(context.TODO())
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
}

func TestAddDirectConnection(t *testing.T) {
	utils.Init()
	middleProxy := CASCADE.Run(true, "http://localhost:8082", "", "")
	var middleServer *http.Server

	go func() {
		utils.Sugar.Info("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	time.Sleep(1 * time.Second)
	_, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)

	Config.LocalPort = "8801"
	AddDirectConnection("google")
	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	ClearHostList()
	_, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)

	err = middleServer.Shutdown(context.TODO())
	assert.NoError(t, err)
}

func TestAddDifferentProxyConnection(t *testing.T) {
	utils.Init()
	username, password := "foo", "bar"
	endProxy := DIRECT.Run(true)
	endProxy.Verbose = true
	var endServer *http.Server
	auth.ProxyBasic(endProxy, "my_realm", func(user, pwd string) bool {
		return user == username && password == pwd
	})
	go func() {
		utils.Sugar.Info("serving end proxy server at localhost:8082")
		endServer = &http.Server{
			Addr:    "localhost:8082",
			Handler: endProxy,
		}
		err := endServer.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	middleProxy := CASCADE.Run(true, "http://localhost:8083", username, password)
	var middleServer *http.Server

	go func() {
		utils.Sugar.Info("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	time.Sleep(1 * time.Second)
	Config.LocalPort = "8801"
	utils.Sugar.Info("starting HTTPS test to fail")
	_, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)

	utils.Sugar.Info("starting HTTP test to fail")
	resp, err := utils.GetResponse("http://localhost:8081", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusInternalServerError)

	AddDifferentProxyConnection("google", "http://localhost:8082")
	utils.Sugar.Info("starting HTTPS test to work")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	utils.Sugar.Info("starting HTTP test")
	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	AddDifferentProxyConnection("google", "")
	utils.Sugar.Info("starting HTTPS test to work")
	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	utils.Sugar.Info("starting HTTP test")
	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	resp, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	utils.Sugar.Info("starting HTTP test")
	resp, err = utils.GetResponse("http://localhost:8081", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	AddDifferentProxyConnection("steam", "")
	assert.NotContains(t, HostList, "")

	ClearHostList()
	_, err = utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)

	err = middleServer.Shutdown(context.TODO())
	assert.NoError(t, err)
	err = endServer.Shutdown(context.TODO())
	assert.NoError(t, err)
}

func TestCascadeProxy_ModeSwitch(t *testing.T) {
	utils.Init()
	middleProxy := CASCADE.Run(true, "http://localhost:8083", "", "")
	var middleServer *http.Server

	go func() {
		utils.Sugar.Info("serving middle proxy server at localhost:8081")
		middleServer = &http.Server{
			Addr:    "localhost:8081",
			Handler: middleProxy,
		}
		err := middleServer.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	time.Sleep(1 * time.Millisecond)
	_, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)

	utils.Sugar.Info("writing to DirectOverrideChan")
	DirectOverrideChan = true
	utils.Sugar.Info("Testing direct Override")
	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = middleServer.Shutdown(context.TODO())
	assert.NoError(t, err)
	DirectOverrideChan = false
}
