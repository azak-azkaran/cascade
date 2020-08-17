package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfFromFile(t *testing.T) {
	fmt.Println("Running: TestGetConfFromFile")
	utils.Init()

	conf, err := GetConfFromFile("./test/test.yml")
	assert.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, "TestHealth", conf.CheckAddress)
	assert.Equal(t, "8888", conf.LocalPort)
	assert.Equal(t, "TestHost", conf.ProxyURL)
	assert.Equal(t, "TestPassword", conf.Password)
	assert.Equal(t, "TestUser", conf.Username)
	assert.Equal(t, int64(5), conf.HealthTime)

	conf, err = GetConfFromFile("noname.yaml")
	assert.Error(t, err)
	assert.Nil(t, conf, "Error could read YAML but should not be able to be")
}

func TestGetConfFromVault(t *testing.T) {
	fmt.Println("Running: TestGetConfFromVault")
	utils.Init()
	StartServer(t, "http://localhost:2000")
	time.Sleep(1 * time.Millisecond)

	_, err := GetConfFromVault("http://localhost:2000", "token", "random")
	assert.Error(t, err)
	sealStatus = false

	conf, err := GetConfFromVault("http://localhost:2000", "token", "random")
	assert.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, "TestHealth", conf.CheckAddress)
	assert.Equal(t, "8888", conf.LocalPort)
	assert.Equal(t, "TestHost", conf.ProxyURL)
	assert.Equal(t, "TestPassword", conf.Password)
	assert.Equal(t, "TestUser", conf.Username)
	assert.Equal(t, int64(30), conf.HealthTime)
	//assert.Equal(t, true, conf.CascadeMode)
	assert.Equal(t, false, conf.DisableAutoChangeMode)
}

func TestCreateConfig(t *testing.T) {
	fmt.Println("Running: TestCreateConfig")
	utils.Init()
	Config := Yaml{LocalPort: "8888", CheckAddress: "https://www.google.de", HealthTime: 5, HostList: "google,eclipse,blub", Log: "info"}
	CreateConfig(&Config)

	conf := GetConfig()
	assert.Equal(t, "https://www.google.de", conf.CheckAddress)
	assert.Equal(t, "8888", conf.LocalPort)
	assert.Equal(t, int64(5), conf.HealthTime)
	assert.Equal(t, "", conf.VaultAddr)
	assert.Equal(t, "", conf.VaultToken)
	assert.Equal(t, len(conf.proxyRedirectList), 3)

	Config.VaultToken = "token"
	Config.VaultAddr = "http://localhost:2000"

	StartServer(t, "http://localhost:2000")
	sealStatus = false
	time.Sleep(1 * time.Millisecond)
	CreateConfig(&Config)

	conf = GetConfig()
	assert.Equal(t, VAULT_TEST_HEALTH, conf.CheckAddress)
	assert.Equal(t, "8888", conf.LocalPort)
	assert.Equal(t, int64(30), conf.HealthTime)
	assert.Equal(t, "http://localhost:2000", conf.VaultAddr)
	assert.Equal(t, "token", conf.VaultToken)
	assert.Equal(t, len(conf.proxyRedirectList), 2)
}

func TestSetConfig(t *testing.T) {
	fmt.Println("Running: TestUpdateConfig")
	utils.Init()
	StartServer(t, "http://localhost:2000")
	time.Sleep(1 * time.Millisecond)
	sealStatus = false

	conf, err := GetConfFromFile("./test/test.yml")
	assert.NoError(t, err)
	require.NotNil(t, conf)
	conf.ConfigFile = ""
	assert.Equal(t, "TestHealth", conf.CheckAddress)
	assert.Equal(t, "8888", conf.LocalPort)
	assert.Equal(t, "TestHost", conf.ProxyURL)
	assert.Equal(t, "TestPassword", conf.Password)
	assert.Equal(t, "TestUser", conf.Username)
	assert.Equal(t, int64(5), conf.HealthTime)
	assert.Equal(t, "http://localhost:2000", conf.VaultAddr)
	assert.Equal(t, "random", conf.VaultToken)

	conf = SetConfig(conf)
	require.NotNil(t, currentConfig)

	assert.Equal(t, int64(5), currentConfig.HealthTime)
	//assert.Equal(t, true, conf.CascadeMode)
	assert.Equal(t, false, currentConfig.DisableAutoChangeMode)
}
