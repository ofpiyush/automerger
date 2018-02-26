package automerger

import (
	"net"
	"net/http"
)

var defaultBranch = "master"

func Serve(config *Config) {
	listener, err := net.Listen("tcp", config.Address)
	fatalErr(err)
	tokener := &Token{
		Config: config,
	}
	pushHandler := &PushEventHandler{
		Config:   config,
		GetToken: tokener.GetorDie,
	}

	fatalErr(http.Serve(listener, pushHandler))

}
