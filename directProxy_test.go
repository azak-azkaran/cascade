package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestDirectProxy_Run(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	time.Sleep(1 * time.Second)
	directProxy := DIRECT.Run(true)
	go func() {
		utils.Init(os.Stdout, os.Stdout, os.Stderr)
		utils.Info.Println("serving end proxy server at localhost:8082")
		err := http.ListenAndServe("localhost:8082", directProxy)
		if err == nil {
			t.Error("Direct Proxy did not shutdown?")
		}
	}()

	utils.Info.Println("waiting for running")
	time.Sleep(1 * time.Second)

	resp, err := utils.GetResponse("http://localhost:8082", "https://www.google.de")
	if err != nil {
		t.Error("Error while client https request resource", err)
	}

	if resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp.Status)
	}

	resp, err = utils.GetResponse("http://localhost:8082", "http://www.google.de")
	if err != nil {
		t.Error("Error while client http request resource", err)
	}

	if resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp.Status)
	}
}
