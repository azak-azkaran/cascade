package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"github.com/azak-azkaran/proxy-go/utils"
	"net/http"
	"os"
	"strings"
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

func run(port string, proxyURL string) {
	utils.Info.Println("Starting Proxy")
	created = true
	server = http.Server{
		Addr: port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			utils.Info.Println("handling Request: ", r.Method)
			if r.Method == http.MethodConnect {
				handleTunneling(w, r, proxyURL)
			} else {
				handleHTTP(w, r, proxyURL)
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
	var password string
	var proxyAddress string
	var user string
	var port string
	flag.StringVar(&password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&proxyAddress, "host", "", "Address of a forward proxy")
	flag.StringVar(&user, "user", "", "Username for authentication to a forward proxy")

	flag.StringVar(&port, "port", ":8888", "Localport on which to run the proxy")
	flag.Parse()

	if len(proxyAddress) > 0 {
		if len(user) > 0 {
			var builder strings.Builder
			builder.WriteString(user)
			builder.WriteString(":")
			builder.WriteString(password)
			builder.WriteString("@")
			builder.WriteString(proxyAddress)
			proxyAddress = builder.String()
		}
		utils.Info.Println("Using ProxyAddress: ", proxyAddress)
	}
	run(port, proxyAddress)
}
