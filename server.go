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

func shutdown(timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err := CURRENT_SERVER.Shutdown(ctx)
	if err != nil {
		utils.Error.Println("Error while shutting down", err)
		return err
	}
	return nil
}

func createServer(proxy *goproxy.ProxyHttpServer, addr string, port string) *http.Server {
	utils.Info.Println("Starting Proxy")
	server := http.Server{
		Addr:    addr + ":" + port,
		Handler: proxy,
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	utils.Info.Println("Starting Listening")
	running = true
	CURRENT_SERVER = &server
	return &server
}
