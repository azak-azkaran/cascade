package main

import (
	"github.com/azak-azkaran/cascade/utils"
	"github.com/azak-azkaran/goproxy"
	"github.com/gin-gonic/gin"
	//"github.com/json-iterator/go"
	"net/http"
)

func ConfigureRouter(proxy *goproxy.ProxyHttpServer, addr string, verbose bool) http.Handler {
	utils.Info.Println("Configurating gin Router")
	if verbose {
		gin.DisableConsoleColor()
	}

	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, Config)
	})
	return r
}
