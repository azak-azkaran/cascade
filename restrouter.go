package main

import (
	"encoding/json"
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
	"net/url"
	"strings"
)

type AddRedirect struct {
	Address string `json:"address"`
	Proxy   string `json:"proxy"`
	Message string `json:"message"`
}

type SetOnlineCheckRequest struct {
	OnlineCheck bool `json:"onlineCheck"`
}

type SetDisableAutoChangeModeRequest struct {
	AutoChangeMode bool `json:"autoChangeMode"`
}

type SetCascadeModeRequest struct {
	CascadeMode bool `json:"cascadeMode"`
}

func ConfigureRouter(proxy *goproxy.ProxyHttpServer, addr string, verbose bool) http.Handler {
	utils.Info.Println("Configurating gin Router")
	if verbose {
		gin.DisableConsoleColor()
	}

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(utils.DefaultLogFormatter))
	r.Use(gin.Recovery())
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, Config)
	})

	r.GET("/getOnlineCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, Config.OnlineCheck)
	})

	r.GET("/getAutoMode", func(c *gin.Context) {
		c.JSON(http.StatusOK, !Config.DisableAutoChangeMode)
	})

	r.POST("/addRedirect", addRedirectFunc)
	r.POST("/setOnlineCheck", setOnlineCheckFunc)
	r.POST("/setAutoMode", setDisableAutoChangeModeFunc)
	r.POST("/setCascadeMode", setCascadeModeFunc)
	return r
}

func setCascadeModeFunc(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var req SetCascadeModeRequest
	err := decoder.Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": html.EscapeString("Problem with Decoding Body"),
		})
	}

	if req.CascadeMode {
		utils.Info.Println("Setting Cascade to: DirectMode")
		ChangeMode(true, true)
	} else {
		utils.Info.Println("Setting Cascade to: CascadeMode")
		ChangeMode(false, true)
	}

	post := gin.H{
		"CascadeMode": Config.CascadeMode,
	}
	c.JSON(http.StatusOK, post)

}

func setDisableAutoChangeModeFunc(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var req SetDisableAutoChangeModeRequest
	err := decoder.Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": html.EscapeString("Problem with Decoding Body"),
		})
	}
	utils.Info.Println("Setting AutoChangeMode to:", req.AutoChangeMode)

	Config.DisableAutoChangeMode = !req.AutoChangeMode

	post := gin.H{
		"AutoChangeMode": !Config.DisableAutoChangeMode,
	}
	c.JSON(http.StatusOK, post)
}

func setOnlineCheckFunc(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var req SetOnlineCheckRequest
	err := decoder.Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": html.EscapeString("Problem with Decoding Body"),
		})
	}
	utils.Info.Println("Setting OnlineCheck to:", req.OnlineCheck)

	Config.OnlineCheck = req.OnlineCheck

	post := gin.H{
		"OnlineCheck": Config.OnlineCheck,
	}
	c.JSON(http.StatusOK, post)
}

func addRedirectFunc(c *gin.Context) {
	decoder := json.NewDecoder(c.Request.Body)
	var req AddRedirect
	err := decoder.Decode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": html.EscapeString("Problem with Decoding Body"),
		})
	}

	utils.Info.Println("Got Address:", req.Address, "\tRedirect to:", req.Proxy, "\n", req)

	proxyURL, err := url.Parse(req.Proxy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": html.EscapeString("Proxy URL could not be parsed"),
		})
	}

	addressURL, err := url.Parse(req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": html.EscapeString("Address URL could not be parsed"),
		})
	}

	if len(Config.HostList) > 0 {
		Config.HostList += ","
	}
	Config.HostList += req.Address + "->" + req.Proxy
	Config.proxyRedirectList = strings.Split(Config.HostList, ",")
	AddDifferentProxyConnection(req.Address, req.Proxy)
	err = SetConf(&Config)
	if err != nil {
		utils.Info.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"address": html.EscapeString(addressURL.String()),
			"proxy":   html.EscapeString(proxyURL.String()),
			"message": html.EscapeString("Added to Redirect List but config file was not updated because:\n" + err.Error()),
		})
	}

	post := gin.H{
		"address": html.EscapeString(addressURL.String()),
		"proxy":   html.EscapeString(proxyURL.String()),
		"message": html.EscapeString("Added to Redirect List, updated File at: " + Config.ConfigFile),
	}
	c.JSON(http.StatusOK, post)
}
