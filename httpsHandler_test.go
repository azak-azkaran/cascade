package main

import (
	"github.com/azak-azkaran/putio-go-aria2/utils"
	"os"
	"reflect"
	"testing"
)

func Test_createDialer(t *testing.T) {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)

	dialer, err := createDialer("")
	if err != nil {
		t.Error("Error while creating Connection")
	}

	if reflect.TypeOf(dialer).String() != "proxy.direct" {
		t.Error("Error Dialer is not of the direct type")
	}

	//go run("8889")
	//for !running {
	//	time.Sleep(1 * time.Second)
	//	utils.Info.Println("waiting for running")
	//}
	//dialer, err = createConnection( "http://localhost:8889")
	//if err != nil {
	//	t.Error("Error while creating Connection")
	//}

	//if reflect.TypeOf(dialer).String() == "proxy.direct" {
	//	t.Error("Error Dialer is of the direct type")
	//}
}
