package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/azak-azkaran/cascade/utils"
)

var version = "undefined"
var closeChan bool
var stopChan = make(chan os.Signal, 2)

func Run(config *Yaml) {
	utils.Sugar.Warn("Creating Configuration")
	config.CascadeMode = true
	config = CreateConfig(config)

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(config)

	lastTime := time.Now()
	utils.Sugar.Info("Starting Selection Process")
	ModeSelection(config)
	utils.Sugar.Info("Starting Running Server")

	RunServer()

	for !closeChan {
		currentDuration := time.Since(lastTime)
		if currentDuration > config.health {
			conf := GetConfig()
			lastTime = time.Now()
			config = CreateConfig(&conf)
			go ModeSelection(config)
			time.Sleep(config.health)
		}
	}

	if closeChan {
		utils.Sugar.Info("Close was set")
		ShutdownCurrentServer()
	}
}

func cleanup() {
	ShutdownCurrentServer()
	closeChan = true

	time.Sleep(1 * time.Second)
	utils.Sugar.Info("Happy Death")
}

func main() {
	utils.Init()

	stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		<-stopChan
		utils.Sugar.Error("Stop was called")
		cleanup()
	}()
	config, err := ParseCommandline()
	if err != nil {
		utils.Sugar.Error("Dying Horribly because problems with Configuration: ", err)
	} else if config != nil {
		Run(config)
	} else {
		utils.Sugar.Info("Version: ", version)
	}
}
