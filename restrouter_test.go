package main

import (
	"encoding/json"
	"fmt"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/urfave/negroni"
	"os"
	"testing"
	"time"
)

func TestRestRouter_GetConfig(t *testing.T) {
	fmt.Println("Running: TestCreateBrokenServer")
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	CreateConfig("8082", "", "", "", "https://www.google.de", 5, "golang.org,youtube.com", "info")
	CreateRestEndpoint("localhost", "8081")

	n := negroni.Classic()
	n.UseFunc(HandleConfig)

	go n.Run(":8081")

	time.Sleep(1 * time.Second)

	client, err := utils.GetClient("", 2)
	if err != nil {
		t.Error(err)
		return
	}

	resp, err := client.Get("http://localhost:8081/config")
	if err != nil {
		t.Error(err)
		return
	}

	if resp == nil || resp.StatusCode != 200 {
		t.Error("Response should be 200", resp.StatusCode)
	}
	if resp.Body != nil {
		fmt.Println(resp.Body)
	}

	decoder := json.NewDecoder(resp.Body)
	var decodedConfig Yaml
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		t.Error(err)
		return
	}
	same := decodedConfig.CascadeMode == Config.CascadeMode &&
		decodedConfig.HealthTime == Config.HealthTime &&
		decodedConfig.CheckAddress == Config.CheckAddress &&
		decodedConfig.HostList == Config.HostList

	if !same {
		t.Error("decodedConfig mismatch with Config")
	}
}
