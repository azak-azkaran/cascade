package main

import (
	"html"
	"net/http"
	"net/url"

	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
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

var error_proxy_parse string = html.EscapeString("Proxy URL could not be parsed")
var error_url_parse string = html.EscapeString("Address URL could not be parsed")
var error_binding string = html.EscapeString("Error while binding JSON: ")

func ConfigureRouter(proxy *goproxy.ProxyHttpServer, addr string, verbose bool) http.Handler {
	utils.Sugar.Info("Configurating gin Router")
	if verbose {
		gin.DisableConsoleColor()
	}

	config := gin.LoggerConfig{
		Formatter: utils.DefaultLogFormatter,
		SkipPaths: []string{"/debug/vars"},
		Output:    utils.GetLogger().Writer(),
	}

	r := gin.New()
	r.Use(gin.LoggerWithConfig(config))
	r.Use(gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/debug/vars", expvar.Handler())

	r.GET("/config", func(c *gin.Context) {

		c.JSON(http.StatusOK, GetConfig())
	})

	r.GET("/getOnlineCheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, GetConfig().OnlineCheck)
	})

	r.GET("/getAutoMode", func(c *gin.Context) {
		c.JSON(http.StatusOK, !GetConfig().DisableAutoChangeMode)
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
		return
	}
	utils.Sugar.Info("Recieved Request: ", req)

	config := GetConfig()
	if req.CascadeMode {
		utils.Sugar.Info("Setting Cascade to: CascadeMode")
		config.CascadeMode = false
		ChangeMode(false, config)
	} else {
		utils.Sugar.Info("Setting Cascade to: DirectMode")
		config.CascadeMode = true
		ChangeMode(true, config)
	}
	conf := CreateConfig(config)
	post := gin.H{
		"CascadeMode": conf.CascadeMode,
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
	utils.Sugar.Info("Recieved Request: ", req, " Setting AutoChangeMode to:", req.AutoChangeMode)
	config := GetConfig()
	config.DisableAutoChangeMode = !req.AutoChangeMode

	conf := SetConfig(config)
	post := gin.H{
		"AutoChangeMode": !conf.DisableAutoChangeMode,
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
	utils.Sugar.Info("Recieved Request: ", req, " Setting OnlineCheck to:", req.OnlineCheck)
	config := GetConfig()
	config.OnlineCheck = req.OnlineCheck

	conf := CreateConfig(config)
	post := gin.H{
		"OnlineCheck": conf.OnlineCheck,
	}

	c.JSON(http.StatusOK, post)
}

func addRedirectFunc(c *gin.Context) {
	var req AddRedirect

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": error_binding + err.Error(),
		})
		return
	}

	utils.Sugar.Info("Recieved Request: ", req,
		"\nGot Address:", req.Address,
		"\nRedirect to:", req.Proxy)

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

	config := GetConfig()
	if len(config.HostList) > 0 {
		config.HostList += ","
	}
	config.HostList += req.Address + "->" + req.Proxy
	AddDifferentProxyConnection(req.Address, req.Proxy)
	SetConfig(config)
	post := gin.H{
		"address": html.EscapeString(addressURL.String()),
		"proxy":   html.EscapeString(proxyURL.String()),
		"message": html.EscapeString("Added to Redirect List, updated File at: " + config.ConfigFile),
	}
	c.JSON(http.StatusOK, post)
}
