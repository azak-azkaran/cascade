package main

import (
	"crypto/tls"
	"github.com/azak-azkaran/proxy-go/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func client(proxyUrl string, requestUrl string) ([]byte, error) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	u, err := url.Parse(proxyUrl)
	if err != nil {
		utils.Error.Println(err)
		return nil, err
	}
	tr := &http.Transport{
		Proxy: http.ProxyURL(u),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(requestUrl)
	if err != nil {
		utils.Error.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		utils.Error.Println(err)
		return nil, err
	}
	return dump, nil
}
