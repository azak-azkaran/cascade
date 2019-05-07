package main

import (
	"flag"
	"github.com/azak-azkaran/cascade/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/signal"
	"time"
)

type conf struct {
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	ProxyURL     string `yaml:"host"`
	LocalPort    string `yaml:"port"`
	CheckAddress string `yaml:"health"`
	HealthTime   int64  `yaml:"health-time"`
}

var CLOSE bool = false

func GetConf(path string) *conf {
	config := conf{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		utils.Error.Printf("yamlFile.Get err   #%v ", err)
		return nil
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		utils.Error.Printf("Unmarshal: %v", err)
		return nil
	}
	return &config
}

func Run(config conf) {
	utils.Info.Println(config)
	utils.Info.Println("Creating Configuration")
	CreateConfig(config.LocalPort, config.ProxyURL, config.Username, config.Password, config.CheckAddress, int(config.HealthTime))
	utils.Info.Println("Starting Proxy with the following flags:")
	utils.Info.Println("Username: ", CONFIG.Username)
	utils.Info.Println("Password: ", CONFIG.Password)
	utils.Info.Println("ProxyUrl: ", CONFIG.ProxyURL)
	utils.Info.Println("Health Address: ", CONFIG.CheckAddress)
	utils.Info.Println("Health Time: ", CONFIG.Health)

	lastTime := time.Now()
	utils.Info.Println("Starting Selection Process and Running Server")
	ModeSelection(CONFIG.CheckAddress)
	for !CLOSE {
		currentDuration := time.Since(lastTime)
		if currentDuration > CONFIG.Health {
			lastTime = time.Now()
			go ModeSelection(CONFIG.CheckAddress)
		}
	}
	if CLOSE {
		ShutdownCurrentServer()
	}
}

func ParseCommandline() *conf {
	config := conf{}
	var configFile string
	flag.StringVar(&config.Password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&config.ProxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&config.Username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&config.LocalPort, "port", "8888", "Port on which to run the proxy")
	flag.StringVar(&config.CheckAddress, "health", "https://www.google.de", "Address which is used for health check if available go to cascade mode")
	flag.Int64Var(&config.HealthTime, "health-time", 5, "Duration between health checks")
	flag.StringVar(&configFile, "config", "", "Ignores other parameters and reads config yaml")
	flag.Parse()

	if len(configFile) > 0 {
		return GetConf(configFile)
	}
	return &config
}

func cleanup() {
	ShutdownCurrentServer()
	utils.Info.Println("Happy Death")
	CLOSE = true
}

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()
	config := ParseCommandline()
	if config != nil {
		Run(*config)
	} else {
		utils.Error.Println("Dying Horribly because problems with Configuration")
	}

}
