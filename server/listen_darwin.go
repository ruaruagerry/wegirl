package server

import (
	"fmt"
	"net/http"
	"wegirl/servercfg"

	"github.com/fvbock/endless"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

// acceptHTTPRequest 监听和接受HTTP
func acceptHTTPRequest() {
	var hh http.Handler
	if servercfg.ForTestOnly {
		// 支持客户端跨域访问
		c := cors.New(cors.Options{
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},           // we need this line for cors to allow cross-origin
			ExposedHeaders:   []string{"Set-Session"}, // we need this line for cors to set Access-Control-Expose-Headers
		})
		hh = c.Handler(rootRouter)
	} else {
		// 对外服务器不应该允许跨域访问
		hh = rootRouter
	}

	portStr := fmt.Sprintf(":%d", servercfg.ServerPort)
	log.Printf("Http server listen at:%d\n", servercfg.ServerPort)

	http.Handle("/", hh)

	err := endless.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatalf("Http server ListenAndServe %d failed:%s\n", servercfg.ServerPort, err)
	}
}
