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

func ModeSelection(config *Yaml) {
	var success bool
	utils.Sugar.Info("Running check on: ", config.CheckAddress)
	rep, err := utils.GetResponse("", config.CheckAddress)
	if err != nil {
		utils.Sugar.Info("Error while checking,", config.CheckAddress, " , ", err)
		success = false
	} else {
		if rep.StatusCode == http.StatusOK {
			utils.Sugar.Debug("Response was: ", rep.Status, "\t", rep.StatusCode)
			success = true
		} else {
			utils.Sugar.Debug("Response was: ", rep.Status)
			success = false
		}
	}

	utils.Sugar.Info("Check returns: ", success)
	if config.CascadeMode {
		utils.Sugar.Info("Current Mode: CascadeMode")
	} else {
		utils.Sugar.Info("Current Mode: DirectMode")
	}

	if !config.DisableAutoChangeMode {
		ChangeMode(success, config)
	} else {
		utils.Sugar.Info("Automatic Change Mode is disabled")
	}
}

func ChangeMode(success bool, config *Yaml) {
	if len(config.ProxyURL) == 0 {
		utils.Sugar.Error("ProxyURL was not set so staying in DirectMode")
		config.CascadeMode = false
		DirectOverrideChan = true
		return
	}

	if (!success && config.CascadeMode && config.OnlineCheck) ||
		(!success && CurrentServer == nil && config.OnlineCheck) {

		// switch to direct mode
		utils.Sugar.Warn("switch to: DirectMode")
		config.CascadeMode = false
		DirectOverrideChan = true
		return

	} else if (success && !config.CascadeMode && config.OnlineCheck) ||
		(success && CurrentServer == nil && config.OnlineCheck) {
		// switch to cascade mode
		utils.Sugar.Warn("switch to: CascadeMode")
		config.CascadeMode = true
		DirectOverrideChan = false
		return
	}

	if (success && config.CascadeMode) || (success && CurrentServer == nil) {
		// switch to direct mode
		utils.Sugar.Warn("switch to: DirectMode")
		config.CascadeMode = false
		DirectOverrideChan = true
		return
		//go Config.DirectFunction()
	} else if (!success && !config.CascadeMode) || (!success && CurrentServer == nil) {

		// switch to cascade mode
		utils.Sugar.Warn("switch to: CascadeMode")
		config.CascadeMode = true
		DirectOverrideChan = false
		return
		//go Config.CascadeFunction()
	}
}
