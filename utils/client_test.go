package utils

import (
	"os"
	"testing"
)

func TestGetResponse(t *testing.T) {
	Init(os.Stdout, os.Stdout, os.Stderr)

	resp, err := GetResponse("", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting without proxy, ", err)
	}
	if resp.StatusCode != 200 {
		t.Error("Google could not be requested, ", resp.Status)
	}
}

func TestGetResponseDump(t *testing.T) {
	Init(os.Stdout, os.Stdout, os.Stderr)

	dump, err := GetResponseDump("", "https://www.google.de")
	if err != nil {
		t.Error("Error while requesting without proxy, ", err)
	}
	if len(dump) == 0 {
		t.Error("Google response was empty")
	}
}
