package main

import (
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
)

type cascadeProxy struct {
}

type hostConfig struct {
	addr      string
	reg       *regexp.Regexp
	proxyAddr string
	regString string
}

var CASCADE cascadeProxy
var DirectOverrideChan bool
var LoginRequired bool
var HostList cmap.ConcurrentMap = cmap.New()
var cascadeMode bool = true
var DialTimeout time.Duration = 5 * time.Second

const (
	httpPrefix      = "http://"
	ProxyAuthHeader = "Proxy-Authorization"
)

func SetBasicAuth(username, password string, req *http.Request) {
	req.Header.Set(ProxyAuthHeader, "Basic "+basicAuth(username, password))
}

func ClearHostList() {
	for content := range HostList.IterBuffered() {
		HostList.Remove(content.Key)
	}
}

func basicAuth(username, password string) string {
	var builder strings.Builder
	builder.WriteString(username)
	builder.WriteString(":")
	builder.WriteString(password)
	return base64.StdEncoding.EncodeToString([]byte(builder.String()))
}

func AddDifferentProxyConnection(host string, proxyAddr string) {
	if host == "" {
		utils.Sugar.Error("Empty host in skiplist found, with redirect to: ", proxyAddr)
		return
	}

	var value hostConfig
	val, ok := HostList.Get(proxyAddr)
	if ok {
		value = val.(hostConfig)
		value.regString += "|"
	}

	value.regString += ".*" + host + ".*"
	value.reg = regexp.MustCompile(value.regString)
	value.addr = host
	if !strings.HasPrefix(proxyAddr, httpPrefix) && len(proxyAddr) > 0 {
		value.proxyAddr = httpPrefix + proxyAddr
	} else {
		value.proxyAddr = proxyAddr
	}
	utils.Sugar.Info("Adding Redirect: ", proxyAddr, " for: ", host, " Regex:", value.regString)
	HostList.Set(proxyAddr, value)
}

func AddDirectConnection(host string) {
	AddDifferentProxyConnection(host, "")
}

func parseProxyUrl(proxyURL string) (*url.URL, error) {
	if strings.HasPrefix(proxyURL, httpPrefix) {
		return url.Parse(proxyURL)
	} else {
		return url.Parse(httpPrefix + proxyURL)
	}
}

func directOverride() bool {
	if DirectOverrideChan {
		cascadeMode = false
	} else {
		cascadeMode = true
	}

	return !cascadeMode
}

func directRedirect(Host string, proxyURL string) (bool, string) {
	for content := range HostList.IterBuffered() {
		val := content.Val.(hostConfig)
		if val.reg.MatchString(Host) {
			utils.Sugar.Warn("Matching Host found: ", Host, " for: ", val.regString)
			utils.Sugar.Debug("Regex: ", val.reg)

			if len(val.proxyAddr) != 0 {
				utils.Sugar.Info("with redirect to: ", val.proxyAddr)
				return false, val.proxyAddr
			} else {
				utils.Sugar.Info("Using direct connection")
				return true, ""
			}
		}
	}
	return false, proxyURL
}

func CustomProxy(proxyURL string) func(req *http.Request) (*url.URL, error) {
	return func(reg *http.Request) (*url.URL, error) {
		if directOverride() {
			return nil, nil
		}

		directly, redirectAddr := directRedirect(reg.Host, proxyURL)
		if directly {
			return nil, nil
		}

		f, err := parseProxyUrl(redirectAddr)
		return f, err
	}
}

func CustomConnectDial(proxyURL string, connectReqHandler func(req *http.Request), server *goproxy.ProxyHttpServer) func(network string, addr string) (net.Conn, error) {
	return func(network string, addr string) (conn net.Conn, e error) {
		if directOverride() {
			return net.DialTimeout(network, addr, DialTimeout)
		}

		directly, redirectAddr := directRedirect(addr, proxyURL)
		if directly {
			return net.DialTimeout(network, addr, DialTimeout)
		}
		f := server.NewConnectDialToProxyWithHandler(redirectAddr, connectReqHandler)
		return f(network, addr)
	}
}

func (cascadeProxy) Run(verbose bool, proxyURL string, username string, password string) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	//proxy.Logger = utils.Info
	proxy.Logger = zap.NewStdLog(utils.Sugar.Desugar())
	proxy.KeepHeader = true

	if !strings.HasPrefix(proxyURL, httpPrefix) {
		proxyURL = httpPrefix + proxyURL
	}

	proxy.Tr.Proxy = CustomProxy(proxyURL)
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

	proxy.OnRequest().Do(goproxy.FuncReqHandler(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		if LoginRequired {
			SetBasicAuth(username, password, req)
		}
		return req, nil
	}))
	return proxy
}
