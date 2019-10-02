package main

import (
	"flag"
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Yaml struct {
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	ProxyURL          string `yaml:"host"`
	LocalPort         string `yaml:"port"`
	CheckAddress      string `yaml:"health"`
	HealthTime        int64  `yaml:"health-time"`
	HostList          string `yaml:"host-list"`
	LogPath           string `yaml:"log-path"`
	proxyRedirectList []string
	health            time.Duration
	verbose           bool
	CascadeMode       bool   `yaml:"CascadeMode"`
	Log               string `yaml:"Log"`
}

var Config Yaml

var version = "undefined"
var closeChan bool
var stopChan = make(chan os.Signal, 2)

// GetConf reads the Configuration from a yaml file at @path
func GetConf(path string) (*Yaml, error) {
	config := Yaml{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %v", err)
	}

	if len(config.LocalPort) == 0 {
		config.LocalPort = "8888"
	}

	if len(config.CheckAddress) == 0 {
		config.CheckAddress = "https://www.google.de"
	}

	if len(config.Log) == 0 {
		config.Log = "WARNING"
	}

	if config.HealthTime == 0 {
		config.HealthTime = 5
	}

	return &config, nil
}

func CreateConfig(localPort string, proxyUrl string, username string, password string, checkAddress string, healthTime int, skipHosts string, logLevel string) {
	Config.LocalPort = localPort
	Config.ProxyURL = proxyUrl
	Config.Username = username
	Config.Password = password

	switch strings.ToUpper(logLevel) {
	case "INFO":
		Config.Log = "INFO"
		Config.verbose = true
		utils.EnableInfo()
		break
	case "ERROR":
		Config.Log = "ERROR"
		Config.verbose = false
		break
	case "WARNING":
	default:
		Config.Log = "WARNING"
		Config.verbose = true
	}

	Config.proxyRedirectList = strings.Split(skipHosts, ",")

	Config.CascadeMode = true
	Config.CheckAddress = checkAddress
	Config.health = time.Duration(healthTime) * time.Second

	utils.Info.Println("Creating Server")
	CurrentServer = CreateServer(Config)
}

func Run(config Yaml) {
	utils.Info.Println(config)
	utils.Info.Println("Creating Configuration")
	CreateConfig(config.LocalPort, config.ProxyURL, config.Username, config.Password, config.CheckAddress, int(config.HealthTime), config.HostList, config.Log)
	utils.Info.Println("Starting Proxy with the following flags:")
	utils.Info.Println("Username: ", Config.Username)
	utils.Info.Println("Password: ", Config.Password)
	utils.Info.Println("ProxyUrl: ", Config.ProxyURL)
	utils.Info.Println("Health Address: ", Config.CheckAddress)
	utils.Info.Println("Health Time: ", Config.health)
	utils.Info.Println("Skip Cascade for Hosts: ", Config.proxyRedirectList)
	utils.Info.Println("Log Level: ", Config.Log)

	lastTime := time.Now()
	utils.Info.Println("Starting Selection Process")
	ModeSelection(Config.CheckAddress)
	utils.Info.Println("Starting Running Server")

	if !config.verbose {
		utils.DisableInfo()
	}

	RunServer()

	for !closeChan {
		currentDuration := time.Since(lastTime)
		if currentDuration > Config.health {
			lastTime = time.Now()
			go ModeSelection(Config.CheckAddress)
			time.Sleep(Config.health)
		}
	}

	if closeChan {
		utils.Info.Println("Close was set")
		ShutdownCurrentServer()
	}
}

func ParseCommandline() (*Yaml, error) {
	config := Yaml{}
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
	flag.StringVar(&config.Log, "log", "WARNING", "Log level INFO, WARNING, ERROR")
	ver := flag.Bool("version", false, "prints out the version")
	flag.Parse()

	if *ver {
		return nil, nil
	}
	if len(configFile) > 0 {
		return GetConf(configFile)
	}
	return &config, nil
}

func cleanup() {
	ShutdownCurrentServer()
	closeChan = true

	time.Sleep(1 * time.Second)
	utils.Info.Println("Happy Death")
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
