package main

import (
	"encoding/json"
	"github.com/azak-azkaran/cascade/utils"
	"net/http"
	"regexp"
)

var configEndpoint *regexp.Regexp

func CreateRestEndpoint(addr string, port string) {
	endpoint := addr + ":" + port + "/config"
	configEndpoint = regexp.MustCompile(endpoint)
}

func HandleConfig(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	utils.Info.Println("Handling Host: ", req.Host, " RequestURI: ", req.RequestURI)
	if configEndpoint.MatchString(req.Host + req.RequestURI) {
		utils.Info.Println("Return Config")
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Config)
	}
	next(w, req)
}
