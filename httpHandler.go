package main

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"io"
	"net/http"
	"net/url"
)

func createTransport(proxyURL string) *http.Transport {
	proxyURI, _ := url.Parse(proxyURL)
	if proxyURL == "" {
		return &http.Transport{
			Proxy: nil,
		}
	} else {
		return &http.Transport{
			Proxy: http.ProxyURL(proxyURI),
		}

	}
}

func handleHTTP(w http.ResponseWriter, req *http.Request, proxyURL string) {
	utils.Info.Println("handle HTTP Request to: ", req.Host)
	resp, err := createTransport(proxyURL).RoundTrip(req)
	if err != nil {
		utils.Error.Println(w, err.Error(), http.StatusServiceUnavailable)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	utils.Info.Println("HTTP Request:")
	utils.Info.Println("Header:\n", resp.Header)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
