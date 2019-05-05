package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"github.com/elazarl/goproxy"
	"net"
	"time"
)

type directProxy struct {
}

var DIRECT = directProxy{}

func (directProxy) run(verbose bool) *goproxy.ProxyHttpServer {
	endProxy := goproxy.NewProxyHttpServer()
	endProxy.Verbose = verbose
	endProxy.Logger = utils.Info
	endProxy.Tr.Proxy = nil
	endProxy.ConnectDial = func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, 5*time.Second)
	}
	return endProxy
}
