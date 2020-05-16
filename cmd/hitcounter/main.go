package main

import (
	"github.com/thebaer/hitcounter"
	"github.com/writeas/web-core/log"
	"os"
)

func main() {
	err := hitcounter.Serve()
	if err != nil {
		log.Error("Unable to serve: %s", err)
		os.Exit(1)
	}
}
