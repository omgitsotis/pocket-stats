package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
		port := app.Int(cli.IntOpt{
			Name:   "port",
			Desc:   "The port to listen on for API GRPC connections",
			Value:  8080,
			EnvVar: "PORT",
		})

		router := mux.NewRouter()
		router.NewRoute().Path("/api/stats").HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hit app"))
			},
		)

		server := http.Server{
			Handler: router,
			Addr:    fmt.Sprintf(":%d", *port),
		}

		log.Infof("pocker stats started with port %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.WithError(err).Fatal("Fatal error while running HTTP server")
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("fatal error in app")
	}
}
