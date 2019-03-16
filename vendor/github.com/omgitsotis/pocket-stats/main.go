package main

import (
	"fmt"
	"net/http"
	"os"

	"database/sql"

	cli "github.com/jawher/mow.cli"
	_ "github.com/lib/pq"
	"github.com/omgitsotis/pocket-stats/pkg/database"
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
		database.Init(logger)
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
			Value:  "http://localhost:8080",
			EnvVar: "CALLBACK_URL",
		})

		dbURL := app.String(cli.StringOpt{
			Name:   "db-url",
			Desc:   "The url used to connect to the database",
			Value:  "postgres://lnzigpirpvbhaa:3ed23bbfdacf34260f0282ec003721c0ffe5d8e07d141b001d4ecd9712f06687@ec2-54-235-68-3.compute-1.amazonaws.com:5432/d8r5er79optitn",
			EnvVar: "DATABASE_URL",
		})

		connStr := fmt.Sprintf("%s?sslmode=require", *dbURL)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		pgClient := database.NewPostgresDB(db)

		// TODO move to env
		p := pocket.New(
			"74935-9d486f66d2999047b61328f3",
			&http.Client{},
		)

		redirect := fmt.Sprintf("%s/api/pocket/auth/received", *callbackURL)
		s := server.New(p, redirect, pgClient)

		r := router.CreateRouter(s)
		server := router.NewServer(r, *port)
		router.StartServer(server)
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("fatal error in app")
	}
}