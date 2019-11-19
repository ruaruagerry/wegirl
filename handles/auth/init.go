package auth

import "wegirl/server"

func init() {
	server.RegisterPostHandleNoUserID("/auth/login", loginHandle)
}
