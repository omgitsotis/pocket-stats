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
		logLvl := app.String(cli.StringOpt{
			Name:   "log-level",
			Desc:   "The log level of the app",
			EnvVar: "LOG_LEVEL",
			Value:  "debug",
		})

		level := convertLogLevel(*logLvl)

		format := new(log.JSONFormatter)
		format.TimestampFormat = "02-01-2006 15:04:05"
		log.SetFormatter(format)
		log.SetLevel(level)

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
			EnvVar: "DATABASE_URL",
		})

		pocketConsumerID := app.String(cli.StringOpt{
			Name:   "pocket-consumer-id",
			Desc:   "The consumer ID for the API",
			EnvVar: "POCKET_CONSUMER_ID",
			Value:  "74935-9d486f66d2999047b61328f3",
		})

		authUser := app.String(cli.StringOpt{
			Name:   "auth-user",
			Desc:   "The authorised user for the app",
			EnvVar: "AUTH_USER",
			Value:  "test-user",
		})

		authPass := app.String(cli.StringOpt{
			Name:   "auth-password",
			Desc:   "The authorised password for the app",
			EnvVar: "AUTH_PASSWORD",
			Value:  "test-password",
		})

		// Create connection to postgres database
		connStr := fmt.Sprintf("%s?sslmode=require", *dbURL)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()

		// Create New Postgres client
		pgClient := database.NewPostgresDB(db)

		// Create new pocket client
		p := pocket.New(*pocketConsumerID, &http.Client{})

		// Create new server
		redirect := fmt.Sprintf("%s/api/pocket/auth/received", *callbackURL)
		s := server.New(p, redirect, pgClient, *authUser, *authPass)

		// Create routes
		r := router.CreateRouter(s)
		server := router.NewServer(r, *port)

		// Start server
		router.StartServer(server)
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("fatal error in app")
	}
}

// convertLogLevel converts the info level to the logrus level
func convertLogLevel(lvlString string) log.Level {
	level, err := log.ParseLevel(lvlString)
	if err != nil {
		log.WithError(err).Panic("Error parsing logLevel")
	}

	return level
}
