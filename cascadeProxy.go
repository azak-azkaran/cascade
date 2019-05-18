package main

import (
	"encoding/base64"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"github.com/orcaman/concurrent-map"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type cascadeProxy struct {
}

var CASCADE cascadeProxy
var LoginRequired bool
var hostList cmap.ConcurrentMap = cmap.New()

const (
	ProxyAuthHeader = "Proxy-Authorization"
)

func SetBasicAuth(username, password string, req *http.Request) {
	req.Header.Set(ProxyAuthHeader, "Basic "+basicAuth(username, password))
}

func ClearHostList(){
	for content := range hostList.IterBuffered(){
		hostList.Remove(content.Key)
	}
}

func basicAuth(username, password string) string {
	var builder strings.Builder
	builder.WriteString(username)
	builder.WriteString(":")
	builder.WriteString(password)
	return base64.StdEncoding.EncodeToString([]byte(builder.String()))
}

func HandleDirectHttpRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response){
	var resp *http.Response
	var err error
	if req.URL.Scheme == "" {
		resp,err = utils.GetResponse("", req.Host)
	}else {
		resp,err = utils.GetResponse("", req.URL.Scheme + "://" + req.Host)
	}
    if err != nil {
    	utils.Error.Println("Problem while trying direct connection: ", err)
    	return req, nil
	}
    return req, resp
    }


func AddDirectConnection(server *goproxy.ProxyHttpServer, host string ){
	reg := regexp.MustCompile(".*"+ host + ".*")
	server.OnRequest(goproxy.ReqHostMatches(reg)).DoFunc(
		HandleDirectHttpRequest)

	hostList.SetIfAbsent(host,reg)
}

func CustomConnectDial(proxyURL string, connectReqHandler func(req *http.Request), server *goproxy.ProxyHttpServer) func(network string, addr string) (net.Conn, error){
    return func(network string, addr string) (conn net.Conn, e error) {

    	for content := range hostList.IterBuffered() {
    		val := content.Val
            if val.(*regexp.Regexp).MatchString(addr) {
				return net.DialTimeout(network, addr, 5*time.Second)
			}
		}

		f  := server.NewConnectDialToProxyWithHandler(proxyURL, connectReqHandler)
		return f(network, addr)
	}
}

func (cascadeProxy) Run(verbose bool, proxyURL string, username string, password string) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	proxy.Logger = utils.Info
	proxy.KeepHeader = true

	proxy.Tr.Proxy = func(req *http.Request) (*url.URL, error) {
		if strings.HasPrefix( proxyURL,"http://") {
			return url.Parse(proxyURL)
		} else {
			return url.Parse("http://" + proxyURL)
		}
	}
	var connectReqHandler func(req *http.Request)

	if len(username) > 0 {
		LoginRequired = true
		connectReqHandler = func(req *http.Request) {
			SetBasicAuth(username, password, req)
		}
	} else {
		LoginRequired = false
		connectReqHandler = nil
	}

	proxy.ConnectDial = CustomConnectDial(proxyURL, connectReqHandler, proxy)
	//func(network string, addr string) (conn net.Conn, e error) {
	//}

	proxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if LoginRequired {
			SetBasicAuth(username, password, req)
		}
		return req, nil
	}))
	return proxy
}
