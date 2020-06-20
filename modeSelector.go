package main

import (
	"net/http"
	"strings"

	"github.com/azak-azkaran/cascade/utils"
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
	utils.Sugar.Info("Running check on: ", checkAddress)
	rep, err := utils.GetResponse("", checkAddress)
	if err != nil {
		utils.Sugar.Info("Error while checking,", checkAddress, " , ", err)
		success = false
	} else {
		if rep.StatusCode == http.StatusOK {
			utils.Sugar.Info("Response was: ", rep.Status, "\t", rep.StatusCode)
			success = true
		} else {
			utils.Sugar.Info("Response was: ", rep.Status)
			success = false
		}
	}

	utils.Sugar.Info("Check returns: ", success)
	if Config.CascadeMode {
		utils.Sugar.Info("Current Mode: CascadeMode")
	} else {
		utils.Sugar.Info("Current Mode: DirectMode")
	}

	if !Config.DisableAutoChangeMode {
		ChangeMode(success, Config.OnlineCheck)
	} else {
		utils.Sugar.Info("Automatic Change Mode is disabled")
	}
}

func ChangeMode(success bool, directCheck bool) {
	if len(Config.ProxyURL) == 0 {
		utils.Sugar.Error("ProxyURL was not set so staying in DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		return
	}

	if (!success && Config.CascadeMode && directCheck) ||
		(!success && CurrentServer == nil && directCheck) {

		// switch to direct mode
		utils.Sugar.Warn("switch to: DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		return

	} else if (success && !Config.CascadeMode && directCheck) ||
		(success && CurrentServer == nil && directCheck) {
		// switch to cascade mode
		utils.Sugar.Warn("switch to: CascadeMode")
		Config.CascadeMode = true
		DirectOverrideChan = false
		return
	}

	if (success && Config.CascadeMode) || (success && CurrentServer == nil) {
		// switch to direct mode
		utils.Sugar.Warn("switch to: DirectMode")
		Config.CascadeMode = false
		DirectOverrideChan = true
		return
		//go Config.DirectFunction()
	} else if (!success && !Config.CascadeMode) || (!success && CurrentServer == nil) {

		// switch to cascade mode
		utils.Sugar.Warn("switch to: CascadeMode")
		Config.CascadeMode = true
		DirectOverrideChan = false
		return
		//go Config.CascadeFunction()
	}
}
