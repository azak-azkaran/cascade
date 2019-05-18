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
	HostList     string `yaml:"host-list"`
	LogPath      string `yaml:"log-path"`
}

var closeChan bool = false
var stopChan = make(chan os.Signal, 2)
var LogFile *os.File

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

	if len(config.LocalPort) == 0 {
		config.LocalPort = "8888"
	}

	if len(config.CheckAddress) == 0 {
		config.CheckAddress = "https://www.google.de"
	}

	if config.HealthTime == 0 {
		config.HealthTime = 5
	}

	return &config
}

func Run(config conf) {
	utils.Info.Println(config)
	utils.Info.Println("Creating Configuration")
	CreateConfig(config.LocalPort, config.ProxyURL, config.Username, config.Password, config.CheckAddress, int(config.HealthTime), config.HostList)
	utils.Info.Println("Starting Proxy with the following flags:")
	utils.Info.Println("Username: ", CONFIG.Username)
	utils.Info.Println("Password: ", CONFIG.Password)
	utils.Info.Println("ProxyUrl: ", CONFIG.ProxyURL)
	utils.Info.Println("Health Address: ", CONFIG.CheckAddress)
	utils.Info.Println("Health Time: ", CONFIG.Health)
	utils.Info.Println("Skip Cascade for Hosts: ", CONFIG.SkipCascadeHosts)

	lastTime := time.Now()
	utils.Info.Println("Starting Selection Process and Running Server")
	ModeSelection(CONFIG.CheckAddress)
	for !closeChan {
		currentDuration := time.Since(lastTime)
		if currentDuration > CONFIG.Health {
			lastTime = time.Now()
			go ModeSelection(CONFIG.CheckAddress)
			time.Sleep(CONFIG.Health)
		}
	}
	if closeChan {
		ShutdownCurrentServer()
	}
}

func SetLogPath(path string) *os.File {
	buffer, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		utils.Error.Println("Error while opining Log file:", err)
		return nil
	}
	utils.Init(buffer, buffer, buffer)
	return buffer
}

func ParseCommandline() *conf {
	config := conf{}
	var configFile string
	flag.StringVar(&config.Password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&config.ProxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&config.Username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&config.LocalPort, "port", "8888", "Port on which to run the proxy")
	flag.StringVar(&config.CheckAddress, "health", "https://www.google.de", "Address which is used for health check if available go to direct mode")
	flag.Int64Var(&config.HealthTime, "health-time", 30, "Duration between health checks")
	flag.StringVar(&config.HostList, "host-list", "", "Comma Separated List of Host for which DirectMode is used in Cascade Mode")
	flag.StringVar(&config.LogPath, "log-path", "", "Path to a file to write Log Messages to")
	flag.StringVar(&configFile, "config", "", "Path to config yaml file. If set all other command line parameters will be ignored")
	flag.Parse()

	if len(configFile) > 0 {
		return GetConf(configFile)
	}
	return &config
}

func cleanup() {
	ShutdownCurrentServer()
	closeChan = true
	if LogFile != nil {
		err := LogFile.Close()
		if err != nil {
			utils.Error.Println("Error while closing LogFile Pointer: ", err)
		}
	}

	time.Sleep(2 * time.Second)
	utils.Info.Println("Happy Death")
}

func main() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		<-stopChan
		cleanup()
		//os.Exit(1)
	}()
	config := ParseCommandline()
	if config != nil {
		Run(*config)
	} else {
		utils.Error.Println("Dying Horribly because problems with Configuration")
	}
}
