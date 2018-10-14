package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
	log "github.com/sirupsen/logrus"
)

const (
	name = "pocket-stats-server"
	desc = "An API to generate stats based on my pocket account"
)

func main() {
	app := cli.App(name, desc)

	app.Action = func() {
		log.Debug("Started pocket stats server")
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("fatal error in app")
	}
}
