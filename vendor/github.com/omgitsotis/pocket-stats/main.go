package main

import (
	"fmt"
	"net/http"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/omgitsotis/pocket-stats/pkg/router"
	"github.com/omgitsotis/pocket-stats/pkg/server"
	log "github.com/sirupsen/logrus"
)

const (
	name = "pocket-stats-server"
	desc = "An API to generate stats based on my pocket account"
)

func main() {
	app := cli.App(name, desc)

	app.Before = func() {
		format := new(log.TextFormatter)
		format.TimestampFormat = "02-01-2006 15:04:05"
		format.FullTimestamp = true
		format.ForceColors = true
		log.SetFormatter(format)
		log.SetLevel(log.DebugLevel)

		logger := log.StandardLogger()
		router.Init(logger)
		server.Init(logger)
		pocket.Init(logger)
	}

	app.Action = func() {
		port := app.Int(cli.IntOpt{
			Name:   "port",
			Desc:   "The port to listen on for API GRPC connections",
			Value:  8080,
			EnvVar: "PORT",
		})

		callbackURL := app.String(cli.StringOpt{
			Name:   "callback_url",
			Desc:   "The url used for pocket to call back on validation",
			Value:  "http://localhost",
			EnvVar: "CALLBACK_URL",
		})

		p := pocket.New(
			"74935-9d486f66d2999047b61328f3",
			&http.Client{},
		)

		redirect := fmt.Sprintf("%s/api/pocket/auth/received", *callbackURL)
		s := server.New(p, redirect)

		r := router.CreateRouter(s)
		server := router.NewServer(r, *port)
		router.StartServer(server)
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("fatal error in app")
	}
}
