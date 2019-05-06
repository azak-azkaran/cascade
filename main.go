package main

import (
	"flag"
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"strings"
	"time"
)

func parseCommandline() {
	// TODO: maybe build Integration tests?
	var username string
	var password string
	var proxyURL string
	var localPort string
	var checkAddress string
	var healthTime int

	flag.StringVar(&password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&proxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&localPort, "port", "8888", "Port on which to run the proxy")
	flag.StringVar(&checkAddress, "health", "https://www.google.de", "Address which is used for health check if available go to cascade mode")
	flag.IntVar(&healthTime, "health-time", 5, "Duration between health checks")
	// TODO maybe add configuration yaml file? with proxy exceptions
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
	utils.Info.Println("Health Address: ", checkAddress)
	utils.Info.Println("Health Time: ", healthTime)

	utils.Info.Println("Creating Configuration")
	CreateConfig(localPort, proxyURL, username, password, checkAddress, healthTime)
}

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	lastTime := time.Now()
	utils.Info.Println("Creating Server")
	ModeSelection(CONFIG.CheckAddress)
	utils.Info.Println("Starting Server")
	RunServer()
	for true {
		currentDuration := time.Since(lastTime)
		if currentDuration > CONFIG.Health {
			lastTime = time.Now()
			go ModeSelection(CONFIG.CheckAddress)
		}
	}
}
