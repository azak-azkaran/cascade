package main

import (
	"context"
	"crypto/tls"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"net/http"
	"time"
)

// CurrentServer the running httpserver
var CurrentServer *http.Server

// true if the server has started
var running = false

// RunServer : Starts the Server which selects the Mode in which it should be running,
func RunServer() {
	counter := 5
	for CurrentServer == nil {
		utils.Warning.Println("Server was not created waiting for a one second")
		time.Sleep(1 * time.Second)
		counter = counter - 1
		if counter <= 0 {
			utils.Error.Println("Server was not created in time, going back to ModeSelection")
			ModeSelection(CONFIG.CheckAddress)
			utils.Info.Println("Resetting counter")
			counter = 5
		}
	}
	go func() {
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
func CreateServer(proxy *goproxy.ProxyHttpServer, addr string, port string) *http.Server {
	utils.Info.Println("Starting Proxy")
	CurrentServer = createServer(proxy, addr, port)
	return CurrentServer
}
