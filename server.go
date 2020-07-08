package main

import (
	"context"
	"crypto/tls"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"

	//"github.com/urfave/negroni"
	"net/http"
	"syscall"

	"time"
)

// CurrentServer the running httpserver
var CurrentServer *http.Server

// true if the server has started
var running = false

// RunServer : Starts the Server which selects the Mode in which it should be running,
func RunServer() {
	go func() {
		counter := 5
		for CurrentServer == nil {
			utils.Sugar.Warn("Server was not created waiting for a one second")
			time.Sleep(1 * time.Second)
			counter = counter - 1
			if counter <= 0 {
				utils.Sugar.Error("Server was not created in time")
				stopChan <- syscall.SIGINT
			}
		}

		running = true
		err := CurrentServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				utils.Sugar.Info("Server was closed")
			} else {
				utils.Sugar.Error("Error while running server: ", err)
			}
		}
		running = false
	}()
	utils.Sugar.Info("Server started")
}

//ShutdownCurrentServer Shuts down the current server with a 1 Second timeout
func ShutdownCurrentServer() {
	if CurrentServer != nil {
		shutdown(1*time.Second, CurrentServer)
		CurrentServer = nil
		ClearHostList()
	}
}

func shutdown(timeout time.Duration, server *http.Server) {
	utils.Sugar.Info("Starting shutdown with Timout: ", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		utils.Sugar.Error("Error while shutdown: ", err)
	}
}

func create(proxy *goproxy.ProxyHttpServer, addr string, config *Yaml) *http.Server {
	return &http.Server{
		Addr:    addr + ":" + config.LocalPort,
		Handler: ConfigureRouter(proxy, addr, config.verbose),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

// CreateServer creates a proxy server on addr:port
func CreateServer(config *Yaml) *http.Server {
	utils.Sugar.Warn("Creating Proxy on: localhost", ":", config.LocalPort)
	proxy := CASCADE.Run(config.verbose, config.ProxyURL, config.Username, config.Password)
	HandleCustomProxies(config.proxyRedirectList)
	CurrentServer = create(proxy, "localhost", config)
	return CurrentServer
}
