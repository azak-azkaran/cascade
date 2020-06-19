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
	utils.Warning.Println("Creating Configuration")
	Config = config
	CreateConfig()

	utils.Info.Println("Creating Server")
	CurrentServer = CreateServer(Config)

	lastTime := time.Now()
	utils.Info.Println("Starting Selection Process")
	ModeSelection(Config.CheckAddress)
	utils.Info.Println("Starting Running Server")

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
		utils.Info.Println("Close was set")
		ShutdownCurrentServer()
	}
}

func cleanup() {
	ShutdownCurrentServer()
	closeChan = true

	time.Sleep(1 * time.Second)
	utils.Info.Println("Happy Death")
	utils.Close()
}

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	stopChan = make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		<-stopChan
		utils.Error.Println("Stop was called")
		cleanup()
	}()
	config, err := ParseCommandline()
	if err != nil {
		utils.Error.Println("Dying Horribly because problems with Configuration: ", err)
	} else if config != nil {
		Run(*config)
	} else {
		utils.Info.Println("Version: ", version)
	}
}
