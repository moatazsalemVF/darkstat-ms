package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moatazsalemVF/darkstat-ms/system"
	"github.com/moatazsalemVF/ms-template/utils"
)

func main() {
	utils.Initialize()
	addControllers()
	utils.Router.Run(utils.Conf.Server.Address + ":" + fmt.Sprint(utils.Conf.Server.Port))
}

func addControllers() {
	utils.Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	utils.Router.GET("/stats/local", func(c *gin.Context) {
		localDarkStat(c.Writer, c.Request)
	})

	utils.Router.GET("/stats/we", func(c *gin.Context) {
		weDarkStat(c.Writer, c.Request)
	})

	utils.Router.GET("/stats/orange", func(c *gin.Context) {
		orangeDarkStat(c.Writer, c.Request)
	})

}

type response struct {
	Hosts    []system.Host `json:"hosts"`
	Sequence int64         `json:"sequence"`
	Time     int64         `json:"time"`
}

var localSequence = 0
var weSequence = 0
var orangeSequence = 0

func localDarkStat(w http.ResponseWriter, r *http.Request) {
	localSequence++
	hosts, _ := system.Publish("http://10.10.10.40:662/hosts/?full=yes&sort=total", "10.10.10", 0)
	var res response
	res.Hosts = hosts
	res.Sequence = int64(localSequence)
	res.Time = getTimeStamp()
	jsonResponse, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func weDarkStat(w http.ResponseWriter, r *http.Request) {
	weSequence++
	hosts, _ := system.Publish("http://192.168.10.1:662/hosts/?full=yes&sort=total", "192.168.10", 1)
	var res response
	res.Hosts = hosts
	res.Sequence = int64(weSequence)
	res.Time = getTimeStamp()
	jsonResponse, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func orangeDarkStat(w http.ResponseWriter, r *http.Request) {
	orangeSequence++
	hosts, _ := system.Publish("http://192.168.20.1/api/monitoring/month_statistics", "192.168.20", 2)
	var res response
	res.Hosts = hosts
	res.Sequence = int64(orangeSequence)
	res.Time = getTimeStamp()
	jsonResponse, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
