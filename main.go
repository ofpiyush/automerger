package main

import (
	"github.com/ofpiyush/automerger/automerger"
)

func main() {
	automerger.Serve(automerger.ConfigureOrDie())
}
