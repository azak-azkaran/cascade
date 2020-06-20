package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	vault "github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v2"
)

const ()

type Yaml struct {
	Username              string `yaml:"username"`
	Password              string `yaml:"password"`
	ProxyURL              string `yaml:"host"`
	LocalPort             string `yaml:"port"`
	CheckAddress          string `yaml:"health"`
	HealthTime            int64  `yaml:"health-time"`
	HostList              string `yaml:"host-list"`
	LogPath               string `yaml:"log-path"`
	proxyRedirectList     []string
	health                time.Duration
	verbose               bool
	CascadeMode           bool   `yaml:"cascadeMode"`
	Log                   string `yaml:"log"`
	OnlineCheck           bool   `yaml:"onlineCheck"`
	ConfigFile            string `yaml:"configFile"`
	DisableAutoChangeMode bool   `yaml:"disableAutoChangeMode"`
	VaultToken            string `yaml:"vaultToken"`
	VaultAddr             string `yaml:"vaultAddr"`
}

func SetConf(config *Yaml) error {
	f, err := os.Create(config.ConfigFile)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	encoder := yaml.NewEncoder(w)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	err = encoder.Close()
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func SealStatus(config *vault.Config) (*vault.SealStatusResponse, error) {
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	sys := client.Sys()
	respones, err := sys.SealStatus()
	if err != nil {
		return nil, err
	}
	return respones, nil
}

func GetSecret(config *vault.Config, token string, path string) (*vault.Secret, error) {
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}
	client.SetToken(token)

	logical := client.Logical()
	secret, err := logical.Read("cascade/data/" + path)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func GetConfFromVault(vaultAddr string, vaultToken string, path string) (*Yaml, error) {
	config := Yaml{}

	//vaultConfig := vault.DefaultConfig()
	vaultConfig := &vault.Config{
		Address: vaultAddr,
	}

	resp, err := SealStatus(vaultConfig)
	if err != nil {
		return nil, err
	}

	if resp.Sealed {
		return nil, errors.New("Vault is sealed")
	}

	secret, err := GetSecret(vaultConfig, vaultToken, path)
	if err != nil {
		return nil, err
	}

	data := secret.Data["data"].(map[string]interface{})
	if len(data) == 0 {
		return nil, errors.New("Data of secret with path: " + path + " is empty")
	}

	config.Username = data["username"].(string)
	config.Password = data["password"].(string)
	config.LocalPort = data["port"].(string)
	config.ProxyURL = data["host"].(string)
	config.HostList = data["host-list"].(string)
	config.CheckAddress = data["health"].(string)

	health, err := strconv.ParseInt(data["health-time"].(string), 10, 0)
	if err != nil {
		return nil, err
	}
	config.HealthTime = health

	//disableAutoChangeMode, err := strconv.ParseBool(data["disableAutoChangeMode"].(string))
	//if err != nil {
	//	return nil, err
	//}

	//cascadeMode, err := strconv.ParseBool(data["cascadeMode"].(string))
	//if err != nil {
	//	return nil, err
	//}

	//config.DisableAutoChangeMode = disableAutoChangeMode
	//config.CascadeMode = cascadeMode

	return &config, nil
}

// GetConf reads the Configuration from a yaml file at @path
func GetConfFromFile(path string) (*Yaml, error) {
	config := Yaml{}
	yamlFile, err := ioutil.ReadFile(path)
	config.ConfigFile = path
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

	config.proxyRedirectList = strings.Split(config.HostList, ",")
	return &config, nil
}

func UpdateConfig(config *Yaml) (*Yaml, error) {
	var err error
	if config.ConfigFile != "" && config.ConfigFile != "config" {
		config, err = GetConfFromFile(config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}

	if config.VaultAddr != "" {
		utils.Sugar.Info("Found Vault server address")

		hostname, err := os.Hostname()
		if err != nil {
			utils.Sugar.Error("Error getting hostname: ", err)
		}

		if len(config.VaultToken) == 0 {
			return nil, errors.New("Vault token is not provided")
		}

		config, err = GetConfFromVault(config.VaultAddr, config.VaultToken, hostname)
		if err != nil {
			return nil, err
		}
	}

	config.proxyRedirectList = strings.Split(config.HostList, ",")
	config.health = time.Duration(int(Config.HealthTime)) * time.Second
	return config, nil
}

func CreateConfig() {
	if conf, err := UpdateConfig(&Config); err == nil {
		Config = *conf
	}

	switch strings.ToUpper(Config.Log) {
	case "DEBUG":
		fmt.Println(Config)
		fmt.Println("Starting Proxy with the following flags:")
		fmt.Println("Username: ", Config.Username)
		fmt.Println("Password: ", Config.Password)
		fmt.Println("ProxyUrl: ", Config.ProxyURL)
		fmt.Println("Health Address: ", Config.CheckAddress)
		fmt.Println("Health Time: ", Config.health)
		fmt.Println("Skip Cascade for Hosts: ", Config.proxyRedirectList)
		fmt.Println("Log Level: ", Config.Log)
		utils.EnableDebug()
		Config.Log = "DEBUG"
		Config.verbose = true
	case "INFO":
		Config.Log = "INFO"
		Config.verbose = true
		utils.EnableInfo()
	case "ERROR":
		Config.Log = "ERROR"
		Config.verbose = false
		utils.EnableError()
	case "WARNING":
		fallthrough
	default:
		Config.Log = "WARNING"
		Config.verbose = true
		utils.EnableWarning()
	}
}

func ParseCommandline() (*Yaml, error) {
	config := Yaml{}
	flag.StringVar(&config.Password, "password", "", "Password for authentication to a forward proxy")
	flag.StringVar(&config.ProxyURL, "host", "", "Address of a forward proxy")
	flag.StringVar(&config.Username, "user", "", "Username for authentication to a forward proxy")
	flag.StringVar(&config.LocalPort, "port", "8888", "Port on which to run the proxy")
	flag.StringVar(&config.CheckAddress, "health", "https://www.google.de", "Address which is used for health check if available go to direct mode")
	flag.Int64Var(&config.HealthTime, "health-time", 30, "Duration between health checks")
	flag.StringVar(&config.HostList, "host-list", "", "Comma Separated List of Host for which DirectMode is used in Cascade Mode")
	flag.StringVar(&config.LogPath, "log-path", "", "Path to a file to write Log Messages to")
	flag.StringVar(&config.ConfigFile, "config", "", "Path to config yaml file. If set all other command line parameters will be ignored")
	flag.StringVar(&config.Log, "log", "WARNING", "Log level INFO, WARNING, ERROR")
	flag.BoolVar(&config.DisableAutoChangeMode, "disableAutoChangeMode", false, "Disable the automatically change of the working Modes")

	flag.StringVar(&config.VaultAddr, "vault-addr", "", "Address to the Vault")
	flag.StringVar(&config.VaultToken, "vault-token", "", "Token to the Vault")

	ver := flag.Bool("version", false, "prints out the version")
	flag.Parse()

	if *ver {
		return nil, nil
	}
	return UpdateConfig(&config)
}
