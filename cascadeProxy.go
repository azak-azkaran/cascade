package main

import (
	"encoding/base64"
	"github.com/elazarl/goproxy"
	"net/http"
	"net/url"
	"strings"
)

type cascadeProxy struct {
}

var CASCADE cascadeProxy
var LOGINREQUEIRED bool

const (
	ProxyAuthHeader = "Proxy-Authorization"
)

func SetBasicAuth(username, password string, req *http.Request) {
	req.Header.Set(ProxyAuthHeader, "Basic "+basicAuth(username, password))
}

func basicAuth(username, password string) string {
	var builder strings.Builder
	builder.WriteString(username)
	builder.WriteString(":")
	builder.WriteString(password)
	return base64.StdEncoding.EncodeToString([]byte(builder.String()))
}

func (cascadeProxy) Run(verbose bool, proxyURL string, username string, password string) *goproxy.ProxyHttpServer {
	middleProxy := goproxy.NewProxyHttpServer()
	middleProxy.Verbose = verbose

	middleProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse("http://" + proxyURL)
	}
	var connectReqHandler func(req *http.Request)

	if len(username) > 0 {
		LOGINREQUEIRED = true
		connectReqHandler = func(req *http.Request) {
			SetBasicAuth(username, password, req)
		}
	} else {
		LOGINREQUEIRED = false
		connectReqHandler = nil
	}

	middleProxy.ConnectDial = middleProxy.NewConnectDialToProxyWithHandler(proxyURL, connectReqHandler)
	middleProxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if LOGINREQUEIRED {
			SetBasicAuth(username, password, req)
		}
		return req, nil
	}))
	return middleProxy
}
