package automerger

import (
	"net"
	"net/http"
)

var defaultBranch = "master"

func Serve(config *Config) {
	listener, err := net.Listen("tcp", config.Address)
	fatalErr(err)
	pushHandler := &PushEventHandler{
		Config: config,
	}

	fatalErr(http.Serve(listener, pushHandler))

}
