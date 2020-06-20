package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
)

func TestDirectProxy_Run(t *testing.T) {
	utils.Init()
	time.Sleep(1 * time.Second)
	directProxy := DIRECT.Run(true)
	var directServer *http.Server
	go func() {
		utils.Sugar.Info("serving end proxy server at localhost:8082")
		directServer = &http.Server{
			Addr:    "localhost:8082",
			Handler: directProxy,
		}
		err := directServer.ListenAndServe()
		if err != http.ErrServerClosed {
			t.Error("Other Error then ServerClose", err)
		}
	}()

	utils.Sugar.Info("waiting for running")
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

	err = directServer.Shutdown(context.TODO())
	if err != nil {
		t.Error("Error while shutting down server")
	}
}
