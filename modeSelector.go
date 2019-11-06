package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"net/http"
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
		utils.Info.Println("Error while checking,", checkAddress, " , ", err)
		success = false
	} else {
		if rep.StatusCode == http.StatusOK {
			utils.Info.Println("Response was: ", rep.Status, "\t", rep.StatusCode)
			success = true
		} else {
			utils.Info.Println("Response was: ", rep.Status)
			success = false
		}
	}

	utils.Info.Println("Check returns: ", success)
	if Config.CascadeMode {
		utils.Info.Println("Current Mode: CascadeMode")
	} else {
		utils.Info.Println("Current Mode: DirectMode")
	}
	ChangeMode(success, Config.OnlineCheck)
}

func ChangeMode(success bool, directCheck bool) {
	if len(Config.ProxyURL) == 0 {
		utils.Error.Println("ProxyURL was not set so staying in DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		return
	}

	if (!success && Config.CascadeMode && directCheck) ||
		(!success && CurrentServer == nil && directCheck) {

		// switch to direct mode
		utils.Warning.Println("switch to: DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		return

	} else if (success && !Config.CascadeMode && directCheck) ||
		(success && CurrentServer == nil && directCheck) {
		// switch to cascade mode
		utils.Warning.Println("switch to: CascadeMode")
		Config.CascadeMode = true
		DirectOverrideChan = false
		return

	}

	if (success && Config.CascadeMode) || (success && CurrentServer == nil) {
		// switch to direct mode
		utils.Warning.Println("switch to: DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		return
		//go Config.DirectFunction()
	} else if (!success && !Config.CascadeMode) || (!success && CurrentServer == nil) {

		// switch to cascade mode
		utils.Warning.Println("switch to: CascadeMode")
		Config.CascadeMode = true
		DirectOverrideChan = false
		return
		//go Config.CascadeFunction()
	}
}
