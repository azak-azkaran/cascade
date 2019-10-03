package main

import (
	"encoding/json"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/julienschmidt/httprouter"
	//"github.com/json-iterator/go"
	"net/http"
)

var RestRouter *httprouter.Router

func CreateRestEndpoint(addr string, port string, useNotFoundHandler bool) {

	RestRouter = httprouter.New()
	RestRouter.GET("/config", HandleConfig)
	if !useNotFoundHandler {
		RestRouter.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
	//endpoint := addr + ":" + port + "/config"
	//configEndpoint = regexp.MustCompile(endpoint)
}

func HandleConfig(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	utils.Info.Println("Handling Host: ", req.Host, " RequestURI: ", req.RequestURI)
	//if configEndpoint.MatchString(req.Host + req.RequestURI) {
	utils.Info.Println("Return Config")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Config)
	//}
}
