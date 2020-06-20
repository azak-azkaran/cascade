package main

import (
	"net"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"go.uber.org/zap"
)

type directProxy struct {
}

var DIRECT = directProxy{}

func (directProxy) Run(verbose bool) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	proxy.Logger = zap.NewStdLog(utils.Sugar.Desugar())
	proxy.Tr.Proxy = nil
	proxy.ConnectDial = func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, 5*time.Second)
	}
	return proxy
}
