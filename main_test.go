package main

import (
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestGetConf(t *testing.T) {
	fmt.Println("Running: TestGetConf")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	conf, _ := GetConf("./test/test.yml")

	if conf.CheckAddress != "TestHealth" {
		t.Error("CheckAddress was not read correctly, was: ", conf.CheckAddress)
	}

	if conf.LocalPort != "8888" {
		t.Error("Port was not read correctly, was:", conf.LocalPort)
	}

	if conf.ProxyURL != "TestHost" {
		t.Error("ProxyURL was not read correctly, was: ", conf.ProxyURL)
	}

	if conf.Password != "TestPassword" {
		t.Error("Password was not read correctly, was: ", conf.Password)
	}

	if conf.Username != "TestUser" {
		t.Error("Username was not read correctly, was: ", conf.Username)
	}

	if conf.HealthTime != int64(5) {
		t.Error("HealthTime was not read correctly, was: ", conf.HealthTime)
	}

	conf, _ = GetConf("noname.yaml")
	if conf != nil {
		t.Error("Error could read YAML but should not be able to be")
	}
}

func TestRun(t *testing.T) {
	fmt.Println("Running: TestRun")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	config := Yaml{}
	config.HealthTime = 5
	config.Username = "foo"
	config.Password = "bar"
	config.ProxyURL = "localhost:8082"
	config.LocalPort = "8081"
	config.CheckAddress = "https://google.de"

	go Run(config)

	time.Sleep(1 * time.Second)
	if CurrentServer == nil {
		t.Error("No Server was created")
	}

	if !DirectOverrideChan {
		t.Error("Direct Override is not active")
	}

	resp, err := utils.GetResponse("http://localhost:8081", "https://www.google.de")
	if err != nil {
		t.Error("Error while client request over cascade", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Error while client https Request, ", resp.Status)
	}

	cleanup()
	time.Sleep(1 * time.Second)
	if running {
		t.Error("Server is still running")
	}
	if CurrentServer != nil {
		t.Error("Server was not created")
	}
}

func TestMain(t *testing.T) {
	fmt.Println("Running: TestMain")
	closeChan = false
	go main()

	time.Sleep(2 * time.Second)
	utils.Info.Println("calling HTTP")
	resp, err := utils.GetResponse("http://localhost:8888", "http://www.google.de")
	if err != nil {
		t.Error("http test failed: ", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("http test failed: ", resp)
	}

	utils.Info.Println("calling HTTPs")
	resp, err = utils.GetResponse("http://localhost:8888", "https://www.google.de")
	if err != nil {
		t.Error("http test failed: ", err)
	}
	if resp == nil || resp.StatusCode != 200 {
		t.Error("http test failed: ", resp)
	}

	utils.Info.Println("Closing")
	stopChan <- syscall.SIGINT
	time.Sleep(2 * time.Second)

	if CurrentServer != nil {
		t.Error("Server was not reset")
	}

	if !closeChan {
		t.Error("Server was not closed")
	}
}
