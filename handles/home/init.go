package home

import (
	"math/rand"
	"time"
	"wegirl/server"
)

func init() {
	droprand = rand.New(rand.NewSource(time.Now().UnixNano()))

	server.RegisterGetHandleNoUserID("/home/tags", tagsHandle)
	server.RegisterPostHandleNoUserID("/home/girls", girlsHandle)
}
