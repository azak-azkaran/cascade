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

var currentConfig *Yaml

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
	Verbose               bool   `yaml:"verbose"`
	CascadeMode           bool   `yaml:"cascadeMode"`
	Log                   string `yaml:"log"`
	OnlineCheck           bool   `yaml:"onlineCheck"`
	ConfigFile            string `yaml:"configFile"`
	DisableAutoChangeMode bool   `yaml:"disableAutoChangeMode"`
	VaultToken            string `yaml:"vaultToken"`
	VaultAddr             string `yaml:"vaultAddr"`
}

func WriteConfig(config *Yaml) error {
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

	if secret == nil || secret.Data["data"] == nil {
		return nil, errors.New("cascade/data/" + path + " was empty")
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

	if data["username"] != nil {
		config.Username = data["username"].(string)
	} else {
		return nil, errors.New("Username is missing in Vault")
	}

	if data["password"] != nil {
		config.Password = data["password"].(string)
	} else {
		return nil, errors.New("Password is missing in Vault")
	}

	if data["host"] != nil {
		config.ProxyURL = data["host"].(string)
	} else {
		return nil, errors.New("Host is missing in Vault")
	}

	if data["port"] != nil {
		config.LocalPort = data["port"].(string)
	} else {
		utils.Sugar.Warn("Port is missing in Vault")
		config.LocalPort = "8888"
	}

	if data["host-list"] != nil {
		config.HostList = data["host-list"].(string)
	} else {
		utils.Sugar.Warn("Host list is missing in Vault")
		config.HostList = ""
	}

	if data["health"] != nil {
		config.CheckAddress = data["health"].(string)
	} else {
		config.CheckAddress = "https://www.google.de"
	}

	if data["log"] != nil {
		config.Log = data["log"].(string)
	} else {
		config.Log = "WARNING"
	}

	if data["health-time"] != nil {
		health, err := strconv.ParseInt(data["health-time"].(string), 10, 0)
		if err != nil {
			return nil, err
		}
		config.HealthTime = health
	} else {
		config.HealthTime = 5
	}

	if data["disableAutoChangeMode"] != nil {
		disableAutoChangeMode, err := strconv.ParseBool(data["disableAutoChangeMode"].(string))
		if err != nil {
			return nil, err
		}
		config.DisableAutoChangeMode = disableAutoChangeMode
	} else {
		config.DisableAutoChangeMode = false
	}

	//if data["cascadeMode"] != nil {
	//	cascadeMode, err := strconv.ParseBool(data["cascadeMode"].(string))
	//	if err != nil {
	//		return nil, err
	//	}
	//	config.CascadeMode = cascadeMode
	//}

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

	return &config, nil
}

func getFileConfig(config *Yaml) (*Yaml, error) {
	if config.ConfigFile != "" && config.ConfigFile != "config" {
		utils.Sugar.Info("Check File Configuration")
		file_config, err := GetConfFromFile(config.ConfigFile)
		if err != nil {
			return config, err
		}
		file_config.DisableAutoChangeMode = config.DisableAutoChangeMode
		file_config.CascadeMode = config.CascadeMode
		file_config.ConfigFile = config.ConfigFile
		return file_config, nil
	}
	return config, nil
}

func getVaultConfig(config *Yaml) (*Yaml, error) {
	if config.VaultAddr != "" {
		utils.Sugar.Info("Found Vault server address")

		hostname, err := os.Hostname()
		if err != nil {
			utils.Sugar.Error("Error getting hostname: ", err)
		}

		if len(config.VaultToken) == 0 {
			return config, errors.New("Vault token is not provided")
		}

		vault_config, err := GetConfFromVault(config.VaultAddr, config.VaultToken, hostname)
		if err != nil {
			return config, err
		}
		vault_config.DisableAutoChangeMode = config.DisableAutoChangeMode
		vault_config.VaultAddr = config.VaultAddr
		vault_config.VaultToken = config.VaultToken
		vault_config.CascadeMode = config.CascadeMode
		return vault_config, nil
	}
	return config, nil
}

func CreateConfig(config *Yaml) *Yaml {
	utils.Sugar.Info("Creating configuration")
	conf, err := getFileConfig(config)
	if err != nil {
		utils.Sugar.Error(err.Error())
	}

	conf, err = getVaultConfig(config)
	if err != nil {
		utils.Sugar.Error(err.Error())
	}
	conf = SetConfig(conf)
	return conf
}

func ParseCommandline() *Yaml {
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
		return nil
	}
	return CreateConfig(&config)
}

func GetConfig() *Yaml {
	config := currentConfig
	return config
}

func SetConfig(conf *Yaml) *Yaml {
	conf.proxyRedirectList = strings.Split(conf.HostList, ",")
	conf.health = time.Duration(int(conf.HealthTime)) * time.Second

	switch strings.ToUpper(conf.Log) {
	case "DEBUG":
		utils.Sugar.Debug(conf)
		utils.Sugar.Debug("Starting Proxy with the following flags:")
		utils.Sugar.Debug("Username: ", conf.Username)
		utils.Sugar.Debug("Password: ", conf.Password)
		utils.Sugar.Debug("ProxyUrl: ", conf.ProxyURL)
		utils.Sugar.Debug("Health Address: ", conf.CheckAddress)
		utils.Sugar.Debug("Health Time: ", conf.health)
		utils.Sugar.Debug("Skip Cascade for Hosts: ", conf.proxyRedirectList)
		utils.Sugar.Debug("Log Level: ", conf.Log)
		utils.Sugar.Debug("OnlineCheck: ", conf.OnlineCheck)
		utils.Sugar.Debug("CascadeMode: ", conf.CascadeMode)
		utils.Sugar.Debug("DisableAutoChangeMode: ", conf.DisableAutoChangeMode)
		utils.EnableDebug()
		conf.Log = "DEBUG"
		conf.Verbose = true
	case "INFO":
		conf.Log = "INFO"
		conf.Verbose = false
		utils.EnableInfo()
	case "ERROR":
		conf.Log = "ERROR"
		conf.Verbose = false
		utils.EnableError()
	case "WARNING":
		fallthrough
	default:
		conf.Log = "WARNING"
		conf.Verbose = false
		utils.EnableWarning()
	}
	currentConfig = conf
	return currentConfig
}
