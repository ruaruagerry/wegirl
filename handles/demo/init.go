package demo

import "wegirl/server"

func init() {
	server.RegisterGetHandle("/demon/hello", helloHandle)
}
