package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"github.com/elazarl/goproxy"
)

type config struct {
	CascadeMode      bool
	Username         string
	Password         string
	ProxyURL         string
	LocalPort        string
	Verbose          bool
	CascadeFunction  func()
	DirectFunction   func()
	ShutdownFunction func()
	CheckAddress     string
}

var CONFIG config

func switchMode(server *goproxy.ProxyHttpServer, mode string) {
	ShutdownCurrentServer()
	CreateServer(server, "localhost", CONFIG.LocalPort)
	utils.Info.Println("Starting Server in: ", mode)
	RunServer()
}

func CreateConfig(localPort string, proxyUrl string, username string, password string, checkAddress string) {
	CONFIG.LocalPort = localPort
	CONFIG.ProxyURL = proxyUrl
	CONFIG.Username = username
	CONFIG.Password = password
	CONFIG.Verbose = true
	CONFIG.DirectFunction = func() {
		switchMode(DIRECT.Run(CONFIG.Verbose), "Direct Mode")
	}
	CONFIG.CascadeFunction = func() {
		switchMode(CASCADE.Run(CONFIG.Verbose, CONFIG.ProxyURL, CONFIG.Username, CONFIG.Password), "Cascade Mode")
	}
	CONFIG.CheckAddress = checkAddress
}

func ModeSelection(checkAddress string) {
	success := false
	utils.Info.Println("Running check on: ", checkAddress)
	rep, err := utils.GetResponse("", checkAddress)
	if err != nil {
		utils.Error.Println("Error while checking,", checkAddress, " , ", err)
		success = false
	} else {
		if rep.StatusCode == 200 {
			success = true
		} else {
			utils.Info.Println("Response was: ", rep.Status)
			success = false
		}
	}

	if CONFIG.Verbose {
		utils.Info.Println("Check returns: ", success)
		if CONFIG.CascadeMode {
			utils.Info.Println("Current Mode: CascadeMode")
		} else {
			utils.Info.Println("Current Mode: DirectMode")
		}
	}
	ChangeMode(success)
}

func ChangeMode(selector bool) {
	if !selector && CONFIG.CascadeMode {
		// switch to direct mode
		utils.Info.Println("switch to: DirectMode")
		CONFIG.CascadeMode = false
		go CONFIG.DirectFunction()
	} else if selector && !CONFIG.CascadeMode {
		// switch to cascade mode
		utils.Info.Println("switch to: CascadeMode")
		CONFIG.CascadeMode = true
		go CONFIG.CascadeFunction()
	}
}
