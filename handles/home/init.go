package home

import "wegirl/server"

func init() {
	server.RegisterGetHandleNoUserID("/home/tags", tagsHandle)
}
