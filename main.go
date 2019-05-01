package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

var running bool
var created bool
var server http.Server

func shutdown(timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), timeout)

	if !running {
		return errors.New("Server is not running")
	}
	err := server.Shutdown(ctx)
	if err != nil {
		utils.Error.Println("Error while shutting down", err)
		return err
	}
	running = false
	return nil
}

func run() {
	utils.Info.Println("Starting Proxy")
	created = true
	server = http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.Info.Println("handling Request: ", r.Method)
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	utils.Info.Println("Starting Listening")
	running = true
	utils.Error.Println(server.ListenAndServe())
}

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	var pemPath string
	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")
	var keyPath string
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")
	var proto string
	flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")
	flag.Parse()

	run()
}
