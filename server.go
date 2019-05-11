package main

import (
	"context"
	"crypto/tls"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"net/http"
	"time"
)

var CURRENT_SERVER *http.Server = nil
var running = false

func RunServer() {
	counter := 5
	for CURRENT_SERVER == nil {
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
		err := CURRENT_SERVER.ListenAndServe()
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

func ShutdownCurrentServer() {
	if CURRENT_SERVER == nil {
		return
	}
	utils.Info.Println("Starting shutdown with Timout: ", 5*time.Second)
	err := shutdown(5 * time.Second)
	if err != nil {
		utils.Error.Println("Error while shutdown: ", err)
	}
	CURRENT_SERVER = nil
}

func shutdown(timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err := CURRENT_SERVER.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func CreateServer(proxy *goproxy.ProxyHttpServer, addr string, port string) *http.Server {
	utils.Info.Println("Starting Proxy")
	server := http.Server{
		Addr:    addr + ":" + port,
		Handler: proxy,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	utils.Info.Println("Starting Listening")
	CURRENT_SERVER = &server
	return &server
}
