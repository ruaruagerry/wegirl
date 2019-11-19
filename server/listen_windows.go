package server

import (
	"fmt"
	"log"
	"net/http"
	"wegirl/servercfg"
	"time"

	"github.com/rs/cors"
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

	http.Handle("/", hh)

	portStr := fmt.Sprintf(":%d", servercfg.ServerPort)
	s := &http.Server{
		Addr:           portStr,
		Handler:        nil,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 10,
	}

	log.Printf("Http server listen at:%d\n", servercfg.ServerPort)

	err := s.ListenAndServe()
	if err != nil {
		log.Fatalf("Http server ListenAndServe %d failed:%s\n", servercfg.ServerPort, err)
	}
}
