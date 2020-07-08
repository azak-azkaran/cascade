package utils

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/azak-azkaran/goproxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResponse(t *testing.T) {
	Init()

	resp, err := GetResponse("", "https://www.google.de")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	proxy := goproxy.NewProxyHttpServer()

	proxy.ConnectDial = func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, 5*time.Second)
	}

	var server *http.Server
	go func() {
		Sugar.Info("serving end proxy server at localhost:7082")
		server = &http.Server{
			Addr:    "localhost:7082",
			Handler: proxy,
		}
		err := server.ListenAndServe()
		assert.EqualError(t, err, http.ErrServerClosed.Error())
	}()

	time.Sleep(1 * time.Second)
	resp, err = GetResponse("http://localhost:7082", "https://www.google.de")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = server.Shutdown(context.TODO())
	assert.NoError(t, err)
}

func TestGetResponseDump(t *testing.T) {
	Init()

	dump, err := GetResponseDump("", "https://www.google.de")
	assert.NoError(t, err)
	assert.NotEmpty(t, dump)
}
