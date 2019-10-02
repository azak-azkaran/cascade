package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"strings"
)

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

	if Config.verbose {
		utils.Info.Println("Check returns: ", success)
		if Config.CascadeMode {
			utils.Info.Println("Current Mode: CascadeMode")
		} else {
			utils.Info.Println("Current Mode: DirectMode")
		}
	}
	ChangeMode(success)
}

func ChangeMode(selector bool) {
	if (selector && Config.CascadeMode) || (selector && CurrentServer == nil) || len(Config.ProxyURL) == 0 {
		if len(Config.ProxyURL) == 0 && !selector {
			utils.Error.Println("ProxyURL was not set so staying in DirectMode")
		}
		// switch to direct mode
		utils.Info.Println("switch to: DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		//go Config.DirectFunction()
	} else if (!selector && !Config.CascadeMode) || (!selector && CurrentServer == nil) {
		// switch to cascade mode
		utils.Info.Println("switch to: CascadeMode")
		Config.CascadeMode = true
		DirectOverrideChan = false
		//go Config.CascadeFunction()
	}
}
