package main

import (
	"github.com/azak-azkaran/proxy-go/utils"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type HTTPSConnection struct {
	Host                  string
	DestinationConnection net.Conn
	ClientConnection      net.Conn
}

func createDialer(proxyURL string) (proxy.Dialer, error) {

	if proxyURL == "" {
		return proxy.Direct, nil
	} else {
		uri, err := url.Parse(proxyURL)
		if err != nil {
			utils.Error.Println("Error while parsing Proxy Server: ", err)
			return nil, err
		}

		dialer, err := proxy.FromURL(uri, TCP)
		if err != nil {
			utils.Error.Println("Error while creating dialer over Proxy Server: ", err)
			return nil, err
		}

		return dialer, nil
	}
}

func handleTunneling(w http.ResponseWriter, req *http.Request, proxyURL string) {
	utils.Info.Println("handle Tunnel Request to: ", req.Host)
	destConnection, err := net.DialTimeout("tcp", req.Host, 10*time.Second)
	//dialer, err := createDialer(proxyURL)
	//if err != nil {
	//	utils.Error.Println(w, err.Error(), http.StatusServiceUnavailable)
	//	http.Error(w, err.Error(), http.StatusServiceUnavailable)
	//	return
	//}

	//destConnection, err := dialer.Dial("tcp", req.Host)
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
