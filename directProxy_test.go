package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	utils.Sugar.Info("waiting for running")
	time.Sleep(1 * time.Second)

	resp, err := utils.GetResponse("http://localhost:8082", "https://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	resp, err = utils.GetResponse("http://localhost:8082", "http://www.google.de")
	assert.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	err = directServer.Shutdown(context.TODO())
	assert.NoError(t, err)
}
