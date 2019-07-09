package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestGetConf(t *testing.T) {
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
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	config := Conf{}
	config.HealthTime = 5
	config.Username = "foo"
	config.Password = "bar"
	config.ProxyURL = "localhost:8082"
	config.LocalPort = "8081"
	config.CheckAddress = "https://google.de"

	go Run(config)
	CurrentServer = nil

	time.Sleep(1 * time.Second)
	if CurrentServer == nil {
		t.Error("No Server was created")
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

func TestSetLogPath(t *testing.T) {
	message := "message"
	path := "testInfoBuffer"

	buffer := SetLogPath(path)
	LogFile = buffer

	utils.Info.Println(message)
	utils.Warning.Println(message)
	utils.Error.Println(message)

	err := LogFile.Close()
	if err != nil {
		t.Errorf("%s could not be closed: %s", path, err)
	}
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("error opening file: %v", err)
	}

	m := string(dat)
	if !strings.Contains(m, message) {
		t.Error("File does not contain message")
	}

	if !strings.Contains(m, "INFO") {
		t.Error("File does not contain INFO message")
	}
	if !strings.Contains(m, "WARNING") {
		t.Error("File does not contain WARNING message")
	}
	if !strings.Contains(m, "ERROR") {
		t.Error("File does not contain ERROR message")
	}

	err = os.Remove(path)
	if err != nil {
		t.Errorf("%s could not be deleted", path)
	}

}

func Test_Main(t *testing.T) {
	go main()

	time.Sleep(2 * time.Second)
	if CurrentServer == nil {
		t.Error("Server was not reset")
	}
	stopChan <- syscall.SIGINT
	time.Sleep(2 * time.Second)

	if CurrentServer != nil {
		t.Error("Server was not reset")
	}

	if !closeChan {
		t.Error("Server was not closed")
	}
}
