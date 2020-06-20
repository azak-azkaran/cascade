package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

const (
	VAULT_TEST_USERNAME                 = "TestUser"
	VAULT_TEST_PASSWORD                 = "TestPassword"
	VAULT_TEST_HOST                     = "TestHost"
	VAULT_TEST_HEALTH                   = "TestHealth"
	VAULT_TEST_HEALTH_TIME              = "30"
	VAULT_TEST_HOST_LIST                = "TestHostList,HostList"
	VAULT_TEST_LOG                      = "INFO"
	VAULT_TEST_PORT                     = "8888"
	VAULT_TEST_DISABLE_AUTO_CHANGE_MODE = "false"
	VAULT_TEST_CASCADE_MODE             = "true"
	VAULT_TEST_LOG_LEVEL                = "info"
)

var server *http.Server

var sealStatus bool = true

func StartServer(t *testing.T, address string) {
	if running {
		log.Println("Server already running")
		return
	}
	gin.SetMode(gin.TestMode)
	server = &http.Server{
		Addr:    strings.TrimPrefix(address, "http://"),
		Handler: createHandler(),
	}
	go func() {
		running = true
		log.Println("Starting MOCK server at: ", server.Addr)
		err := server.ListenAndServe()
		require.Equal(t, http.ErrServerClosed, err)
		running = false
	}()
	time.Sleep(1 * time.Millisecond)
	t.Cleanup(StopServer)
}

func StopServer() {
	log.Println("stopping server")
	if running {
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
	time.Sleep(1 * time.Millisecond)
}

func test_cascade(c *gin.Context) {
	log.Println("MOCK-Server: called cascade")
	var msg vault.Secret
	data := make(map[string]interface{})
	secret := make(map[string]string)
	secret["username"] = VAULT_TEST_USERNAME
	secret["password"] = VAULT_TEST_PASSWORD
	secret["port"] = VAULT_TEST_PORT
	secret["health-time"] = VAULT_TEST_HEALTH_TIME
	secret["health"] = VAULT_TEST_HEALTH
	secret["host"] = VAULT_TEST_HOST
	secret["host-list"] = VAULT_TEST_HOST_LIST
	secret["disableAutoChangeMode"] = VAULT_TEST_DISABLE_AUTO_CHANGE_MODE
	secret["cascadeMode"] = VAULT_TEST_CASCADE_MODE
	secret["log"] = VAULT_TEST_LOG_LEVEL
	data["data"] = secret
	msg.Data = data
	c.JSON(http.StatusOK, msg)
}

func createHandler() http.Handler {
	r := gin.Default()
	r.GET("/v1/sys/seal-status", test_seal_status)
	r.GET("/v1/cascade/data/:name", test_cascade)
	return r
}

func test_seal_status(c *gin.Context) {
	log.Println("MOCK-Server: called seal-status")
	var msg vault.SealStatusResponse

	msg.Sealed = sealStatus
	c.JSON(http.StatusOK, msg)

}
