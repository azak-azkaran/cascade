package main

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"net/http"
	"os"
	"testing"
)

func Test_createTransport(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	tr := createTransport("")

	if tr.ProxyConnectHeader != nil {
		t.Error("Proxy was not set correctly")
	}

	if tr.Proxy != nil {
		t.Error("Proxy was not set correctly")
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://www.google.com")

	if err != nil {
		t.Error("error while get: ", err)
	}

	if resp.StatusCode != 200 {
		t.Error("Return code was not ok (200)")
	}

	tr = createTransport("http://localhost:8889")

	if tr.ProxyConnectHeader != nil {
		t.Error("Proxy was not set correctly")
	}

	if tr.Proxy == nil {
		t.Error("Proxy was not set correctly")
	}

	client = &http.Client{Transport: tr}
	resp, err = client.Get("https://www.google.com")
	if err == nil {
		t.Error("Localhost proxy is already running", err)
	}
}
