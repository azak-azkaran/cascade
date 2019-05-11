package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"net"
	"time"
)

type directProxy struct {
}

var DIRECT = directProxy{}

func (directProxy) Run(verbose bool) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	proxy.Logger = utils.Info
	proxy.Tr.Proxy = nil
	proxy.ConnectDial = func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, 5*time.Second)
	}
	return proxy
}
