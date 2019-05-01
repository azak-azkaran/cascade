package main

import (
	"net"
	"time"
)

type tcp struct{}

var TCP = tcp{}

func (tcp) Dial(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, 10*time.Second)
}
