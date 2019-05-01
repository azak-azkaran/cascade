package main

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func createConnection(proxyURL string, host string) (net.Conn, error) {

	if proxyURL == "" {
		return proxy.Direct.Dial("tcp", host)
	} else {
		uri, err := url.Parse(proxyURL)
		if err != nil {
			utils.Error.Println("Error while parsing Proxy Server")
			return nil, err
		}

		dialer, err := proxy.FromURL(uri, proxy.Direct)
		if err != nil {
			utils.Error.Println("Error while creating dialer over Proxy Server")
			return nil, err
		}

		return dialer.Dial("tcp", host)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	utils.Info.Println("handle Tunnel Request")
	destConnection, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		utils.Error.Println(w, err.Error(), http.StatusServiceUnavailable)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		utils.Error.Println(w, "Hijacking not supported", http.StatusInternalServerError)
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConnection, _, err := hijacker.Hijack()
	if err != nil {
		utils.Error.Println(w, err.Error(), http.StatusServiceUnavailable)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConnection, clientConnection)
	go transfer(clientConnection, destConnection)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
