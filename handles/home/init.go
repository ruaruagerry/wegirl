package home

import "wegirl/server"

func init() {
	server.RegisterGetHandleNoUserID("/home/tags", tagsHandle)
	server.RegisterPostHandleNoUserID("/home/girls", girlsHandle)
}
