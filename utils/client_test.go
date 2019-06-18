package utils

import (
	"context"
	"github.com/azak-azkaran/goproxy"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestGetResponse(t *testing.T) {
	Init(os.Stdout, os.Stdout, os.Stderr)

	resp, err := GetResponse("", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting without proxy, ", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Google could not be requested, ", resp.Status)
	}

	proxy := goproxy.NewProxyHttpServer()

	proxy.ConnectDial = func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, 5*time.Second)
	}

	var server *http.Server
	go func() {
		Init(os.Stdout, os.Stdout, os.Stderr)
		Info.Println("serving end proxy server at localhost:8082")
		server = &http.Server{
			Addr:    "localhost:8082",
			Handler: proxy,
		}
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error("Other Error then ServerClose", err)
		}
	}()

	time.Sleep(1 * time.Second)
	resp, err = GetResponse("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting without proxy, ", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Google could not be requested, ", resp.Status)
	}

	err = server.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down server")
	}
}

func TestGetResponseDump(t *testing.T) {
	Init(os.Stdout, os.Stdout, os.Stderr)

	dump, err := GetResponseDump("", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting without proxy, ", err)
	}
	if len(dump) == 0 {
		t.Error("Google response was empty")
	}
}
