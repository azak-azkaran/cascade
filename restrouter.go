package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"github.com/gin-contrib/expvar"
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

var error_decode string = html.EscapeString("Problem with Decoding Body")
var error_proxy_parse string = html.EscapeString("Proxy URL could not be parsed")
var error_url_parse string = html.EscapeString("Address URL could not be parsed")
var error_binding string = html.EscapeString("Error while binding JSON: ")

func ConfigureRouter(proxy *goproxy.ProxyHttpServer, addr string, verbose bool) http.Handler {
	utils.Info.Println("Configurating gin Router")
	if verbose {
		gin.DisableConsoleColor()
	}

	config := gin.LoggerConfig{
		Formatter: utils.DefaultLogFormatter,
		SkipPaths: []string{"/debug/vars"},
	}

	r := gin.New()
	r.Use(gin.LoggerWithConfig(config))
	r.Use(gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/debug/vars", expvar.Handler())

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
	defer c.Request.Body.Close()
	var req SetCascadeModeRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_binding + err.Error(),
		})
	}

	if req.CascadeMode {
		utils.Info.Println("Setting Cascade to: CascadeMode")
		ChangeMode(true, true)
	} else {
		utils.Info.Println("Setting Cascade to: DirectMode")
		ChangeMode(false, true)
	}

	post := gin.H{
		"CascadeMode": Config.CascadeMode,
	}
	c.JSON(http.StatusOK, post)

}

func setDisableAutoChangeModeFunc(c *gin.Context) {
	defer c.Request.Body.Close()
	var req SetDisableAutoChangeModeRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_binding + err.Error(),
		})
		return
	}
	utils.Info.Println("Setting AutoChangeMode to:", req.AutoChangeMode)

	Config.DisableAutoChangeMode = !req.AutoChangeMode

	post := gin.H{
		"AutoChangeMode": !Config.DisableAutoChangeMode,
	}
	c.JSON(http.StatusOK, post)
}

func setOnlineCheckFunc(c *gin.Context) {
	defer c.Request.Body.Close()
	var req SetOnlineCheckRequest

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_binding + err.Error(),
		})
		return
	}

	utils.Info.Println("Setting OnlineCheck to:", req.OnlineCheck)

	Config.OnlineCheck = req.OnlineCheck

	post := gin.H{
		"OnlineCheck": Config.OnlineCheck,
	}
	c.JSON(http.StatusOK, post)
}

func addRedirectFunc(c *gin.Context) {
	defer c.Request.Body.Close()
	var req AddRedirect

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_binding + err.Error(),
		})
	}

	utils.Info.Println("Got Address:", req.Address, "\tRedirect to:", req.Proxy, "\n", req)

	proxyURL, err := url.Parse(req.Proxy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_proxy_parse + err.Error(),
		})
		return
	}

	addressURL, err := url.Parse(req.Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_url_parse + err.Error(),
		})
		return
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
		return
	}

	post := gin.H{
		"address": html.EscapeString(addressURL.String()),
		"proxy":   html.EscapeString(proxyURL.String()),
		"message": html.EscapeString("Added to Redirect List, updated File at: " + Config.ConfigFile),
	}
	c.JSON(http.StatusOK, post)
}
