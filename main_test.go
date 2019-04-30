package main

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"os"
	"testing"
)

func Test_handleHTTP(t *testing.T){
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	go run("","", "http")

	for !running  {
		utils.Info.Println("waiting for running")
	}
	dump := client("http://localhost:8888")

	if len(dump) == 0 {
		t.Error("No dump recieved")
	}
}

func Test_handleTunneling(t *testing.T){
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
	go run("","", "https")
	for !running  {
		utils.Info.Println("waiting for running")
	}
	dump := client("https://localhost:8888")

	if len(dump) == 0 {
		t.Error("No dump recieved")
	}
}