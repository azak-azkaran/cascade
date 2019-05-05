package main

import (
	"context"
	"crypto/tls"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/elazarl/goproxy"
	"net/http"
	"time"
)

var CURRENT_SERVER *http.Server
var running = false

func RunServer() {
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
	utils.Info.Println("Starting shutdown with Timout: ", 5*time.Second)
	err := shutdown(5 * time.Second)
	if err != nil {
		utils.Error.Println("Error while shutdown: ", err)
	}
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
