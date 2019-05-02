package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	ProxyAuthHeader = "Proxy-Authorization"
)

func SetBasicAuth(username, password string, req *http.Request) {
	req.Header.Set(ProxyAuthHeader, fmt.Sprintf("Basic %s", basicAuth(username, password)))
}

func basicAuth(username, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}

func GetBasicAuth(req *http.Request) (username, password string, ok bool) {
	auth := req.Header.Get(ProxyAuthHeader)
	if auth == "" {
		return
	}

	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

func main() {
	var username string
	var password string
	var proxyURL string
	var localport string

	flag.StringVar(&password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&proxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&localport, "port", "8888", "Localport on which to run the proxy")
	flag.Parse()

	var builder strings.Builder
	builder.WriteString("localhost:")
	builder.WriteString(localport)

	localaddress := builder.String()

	log.Println("Starting Proxy with the following flags:")
	log.Println("Username: ", username)
	log.Println("Password: ", password)
	log.Println("Localaddress: ", localaddress)

	middleProxy := goproxy.NewProxyHttpServer()
	middleProxy.Verbose = true
	middleProxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		return url.Parse(localaddress)
	}
	connectReqHandler := func(req *http.Request) {
		SetBasicAuth(username, password, req)
	}
	middleProxy.ConnectDial = middleProxy.NewConnectDialToProxyWithHandler(proxyURL, connectReqHandler)
	middleProxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		SetBasicAuth(username, password, req)
		return req, nil
	}))
	log.Println("serving middle proxy server at ", localaddress)
	log.Println(http.ListenAndServe(localaddress, middleProxy))
}
