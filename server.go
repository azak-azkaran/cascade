package main

import (
	"context"
	"crypto/tls"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
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
			utils.Warning.Println("Server was not created waiting for a one second")
			time.Sleep(1 * time.Second)
			counter = counter - 1
			if counter <= 0 {
				utils.Error.Println("Server was not created in time")
				stopChan <- syscall.SIGINT
			}
		}

		running = true
		err := CurrentServer.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				utils.Info.Println("Server was closed")
			} else {
				utils.Error.Println("Error while running server: ", err)
			}
		}
		running = false
	}()
	utils.Info.Println("Server started")
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
	utils.Info.Println("Starting shutdown with Timout: ", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		utils.Error.Println("Error while shutdown: ", err)
	}
}

func createServer(proxy *goproxy.ProxyHttpServer, addr string, port string) *http.Server {
	return &http.Server{
		Addr:    addr + ":" + port,
		Handler: proxy,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

// CreateServer creates a proxy server on addr:port
func CreateServer(config config) *http.Server {
	utils.Info.Println("Creating Proxy on: localhost", ":", config.LocalPort)
	proxy := CASCADE.Run(config.Verbose, config.ProxyURL, config.Username, config.Password)
	HandleCustomProxies(config.ProxyRedirectList)
	CurrentServer = createServer(proxy, "localhost", config.LocalPort)
	return CurrentServer
}
