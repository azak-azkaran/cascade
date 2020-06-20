package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/azak-azkaran/cascade/utils"
)

var Config Yaml

var version = "undefined"
var closeChan bool
var stopChan = make(chan os.Signal, 2)

func Run(config Yaml) {
	utils.Sugar.Warn("Creating Configuration")
	Config = config
	Config.CascadeMode = true
	CreateConfig()

	utils.Sugar.Info("Creating Server")
	CurrentServer = CreateServer(Config)

	lastTime := time.Now()
	utils.Sugar.Info("Starting Selection Process")
	ModeSelection(Config.CheckAddress)
	utils.Sugar.Info("Starting Running Server")

	RunServer()

	for !closeChan {
		currentDuration := time.Since(lastTime)
		if currentDuration > Config.health {
			lastTime = time.Now()
			CreateConfig()
			go ModeSelection(Config.CheckAddress)
			time.Sleep(Config.health)
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
		Run(*config)
	} else {
		utils.Sugar.Info("Version: ", version)
	}
}
