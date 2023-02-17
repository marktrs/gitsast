package main

import (
	"os"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/cmd/server"
)

func main() {
	if err := server.Start(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
