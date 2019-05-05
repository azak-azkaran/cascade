package main

import (
	"flag"
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"strings"
	"time"
)

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	var username string
	var password string
	var proxyURL string
	var localPort string
	var checkAddress string

	flag.StringVar(&password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&proxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&localPort, "port", "8888", "Port on which to run the proxy")
	flag.StringVar(&checkAddress, "health", "https://www.google.de", "Address which is used for health check if available go to cascade mode")
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

	CreateConfig(localPort, proxyURL, username, password, checkAddress)

	lastTime := time.Now()
	RunServer()
	for true {
		currentDuration := time.Since(lastTime)
		if currentDuration > 2*time.Second {
			lastTime = time.Now()
			go ModeSelection(CONFIG.CheckAddress)
		}
	}
}
