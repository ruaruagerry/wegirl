// Package server HTTP 服务器
package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/julienschmidt/httprouter"
)

const (
	versionCode = 3.0
)

var (
	// 根router，只有http server看到
	rootRouter = httprouter.New()

	// rootPath = "/game/:uuid/goddess"
	rootPath = ""
)

// GetVersion 版本号
func GetVersion() int {
	return versionCode
}

func echoVersion(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Write([]byte(fmt.Sprintf("version:%v", versionCode)))
}

// CreateHTTPServer 启动服务器
func CreateHTTPServer() {
	log.Printf("CreateHTTPServer")

	rootRouter.Handle("GET", rootPath+"/version", echoVersion)

	redisStartup()
	mysqlStartup()
	// mqStartup()

	go acceptHTTPRequest()
}
