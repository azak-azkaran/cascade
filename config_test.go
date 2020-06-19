package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfFromFile(t *testing.T) {
	fmt.Println("Running: TestGetConfFromFile")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

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
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
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
	assert.Equal(t, true, conf.CascadeMode)
	assert.Equal(t, false, conf.DisableAutoChangeMode)
}

func TestUpdateConfig(t *testing.T) {
	fmt.Println("Running: TestUpdateConfig")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
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
	conf, err = UpdateConfig(*conf)
	assert.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, int64(30), conf.HealthTime)
	assert.Equal(t, true, conf.CascadeMode)
	assert.Equal(t, false, conf.DisableAutoChangeMode)
}
