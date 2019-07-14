package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"strings"
	"time"
)

type config struct {
	CascadeMode       bool   `json:"CascadeMode"`
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	ProxyURL          string `json:"ProxyURL"`
	LocalPort         string
	Verbose           bool
	ShutdownFunction  func()
	CheckAddress      string `json:"CheckAddress"`
	Health            time.Duration
	ProxyRedirectList []string `json:"host-list"`
}

var CONFIG config

func CreateConfig(localPort string, proxyUrl string, username string, password string, checkAddress string, healthTime int, skipHosts string) {
	CONFIG.LocalPort = localPort
	CONFIG.ProxyURL = proxyUrl
	CONFIG.Username = username
	CONFIG.Password = password
	CONFIG.Verbose = true
	CONFIG.ProxyRedirectList = strings.Split(skipHosts, ",")

	CONFIG.CascadeMode = true
	CONFIG.CheckAddress = checkAddress
	CONFIG.Health = time.Duration(healthTime) * time.Second

	utils.Info.Println("Creating Server")
	//switchMode(server, "Cascade Mode")
	CurrentServer = CreateServer(CONFIG)
}

func HandleCustomProxies(list []string) {
	if len(list) == 0 {
		return
	}

	for i := 0; i < len(list); i++ {
		rule := list[i]
		if strings.Contains(rule, "->") {
			split := strings.Split(rule, "->")
			AddDifferentProxyConnection(split[0], split[1])
		} else {
			AddDirectConnection(rule)
		}
	}
}

func ModeSelection(checkAddress string) {
	var success bool
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
	if (selector && CONFIG.CascadeMode) || (selector && CurrentServer == nil) || len(CONFIG.ProxyURL) == 0 {
		if len(CONFIG.ProxyURL) == 0 && !selector {
			utils.Error.Println("ProxyURL was not set so staying in DirectMode")
		}
		// switch to direct mode
		utils.Info.Println("switch to: DirectMode")
		CONFIG.CascadeMode = false
		DirectOverrideChan = true
		//go CONFIG.DirectFunction()
	} else if (!selector && !CONFIG.CascadeMode) || (!selector && CurrentServer == nil) {
		// switch to cascade mode
		utils.Info.Println("switch to: CascadeMode")
		CONFIG.CascadeMode = true
		DirectOverrideChan = false
		//go CONFIG.CascadeFunction()
	}
}
