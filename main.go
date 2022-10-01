package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moatazsalemVF/ms-template/utils"
)

func main() {
	utils.Initialize()
	addPingController()
	utils.Router.Run(utils.Conf.Server.Address + ":" + fmt.Sprint(utils.Conf.Server.Port))
}

func addPingController() {
	utils.Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
