package utils

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func GetResponse(proxyUrl string, requestUrl string) (*http.Response, error) {
	return getResponse(proxyUrl, requestUrl, true)
}

func GetClient(proxyUrl string, timeout int) (*http.Client, error) {
	var tr *http.Transport
	if len(proxyUrl) > 0 {
		u, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		tr = &http.Transport{
			Proxy: http.ProxyURL(u),
			// Disable HTTP/2.
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		}
	} else {
		tr = &http.Transport{
			Proxy: nil,
			// Disable HTTP/2.
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		}
	}

	return &http.Client{Transport: tr, Timeout: time.Duration(time.Duration(timeout) * time.Second)}, nil
}

func getResponse(proxyUrl string, requestUrl string, close bool) (*http.Response, error) {
	client, _ := GetClient(proxyUrl, 2)
	resp, err := client.Get(requestUrl)
	if err != nil {
		return resp, err
	}
	if close {
		defer resp.Body.Close()
	}
	return resp, nil
}

func GetResponseDump(proxyUrl string, requestUrl string) ([]byte, error) {
	resp, err := getResponse(proxyUrl, requestUrl, false)
	if err != nil {
		return nil, err
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return dump, nil
}
