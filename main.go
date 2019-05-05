package main

import (
	"flag"
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"strings"
)

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	var username string
	var password string
	var proxyURL string
	var localPort string

	flag.StringVar(&password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&proxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&localPort, "port", "8888", "Localport on which to run the proxy")
	flag.Parse()

	var builder strings.Builder
	builder.WriteString("localhost:")
	builder.WriteString(localPort)

	localAddress := builder.String()

	utils.Info.Println("Starting Proxy with the following flags:")
	utils.Info.Println("Username: ", username)
	utils.Info.Println("Password: ", password)
	utils.Info.Println("ProxyUrl: ", proxyURL)
	utils.Info.Println("Local Address: ", localAddress)

	utils.Warning.Println("serving middle proxy server at ", localAddress)
	middleProxy := CASCADE.run(true, proxyURL, username, password)
	proxyServer := createServer(middleProxy, "localhost", localPort)
	utils.Error.Println(proxyServer.ListenAndServe())
}
