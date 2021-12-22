package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/omgitsotis/pocket-stats/pkg/database"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/omgitsotis/pocket-stats/pkg/router"
	"github.com/omgitsotis/pocket-stats/pkg/server"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const (
	appName = "pocket-stats-server"

	flagLogLevel         = "log-level"
	flagPort             = "port"
	flagCallbackURL      = "callback_url"
	flagDBURL            = "db-url"
	flagPocketConsumerID = "pocket-consumer-id"
	flagAuthUser         = "auth-user"
	flagAuthPassword     = "auth-password"
)

func main() {
	app := &cli.App{
		Name: appName,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagLogLevel,
				Usage:   "The log level of the app",
				EnvVars: []string{"LOG_LEVEL"},
				Value:   "debug",
			},
			&cli.Int64Flag{
				Name:    flagPort,
				Usage:   "The port to listen on for API GRPC connections",
				Value:   8080,
				EnvVars: []string{"PORT"},
			},
			&cli.StringFlag{
				Name:    flagCallbackURL,
				Usage:   "The url used for pocket to call back on validation",
				Value:   "http://localhost:8080",
				EnvVars: []string{"CALLBACK_URL"},
			},
			&cli.StringFlag{
				Name:    flagDBURL,
				Usage:   "The url used to connect to the database",
				EnvVars: []string{"DATABASE_URL"},
			},
			&cli.StringFlag{
				Name:    flagPocketConsumerID,
				Usage:   "The consumer ID for the API",
				EnvVars: []string{"POCKET_CONSUMER_ID"},
				Value:   "74935-9d486f66d2999047b61328f3",
			},
			&cli.StringFlag{
				Name:    flagAuthUser,
				Usage:   "The authorised user for the app",
				EnvVars: []string{"AUTH_USER"},
				Value:   "test-user",
			},
			&cli.StringFlag{
				Name:    flagAuthPassword,
				Usage:   "The authorised password for the app",
				EnvVars: []string{"AUTH_PASSWORD"},
				Value:   "test-password",
			},
		},
		Before: func(c *cli.Context) error {
			level := convertLogLevel(c.String(flagLogLevel))

			format := new(log.JSONFormatter)
			format.TimestampFormat = "02-01-2006 15:04:05"
			log.SetFormatter(format)
			log.SetLevel(level)

			logger := log.StandardLogger()
			router.Init(logger)
			server.Init(logger)
			pocket.Init(logger)
			return nil
		},
		Action: func(c *cli.Context) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Create connection to postgres database
			connStr := fmt.Sprintf("%s?sslmode=require", c.String(flagDBURL))
			conn := database.MustConnect(ctx, connStr)
			defer conn.Close()

			// Create New Postgres client
			pgClient := database.NewStore(conn)

			// Create new pocket client
			p := pocket.New(c.String(flagPocketConsumerID), &http.Client{})

			// Create new server
			redirect := fmt.Sprintf("%s/api/pocket/auth/received", c.String(flagCallbackURL))
			s := server.New(ctx, p, redirect, pgClient, c.String(flagAuthUser), c.String(flagAuthPassword))

			// Create routes
			r := router.CreateRouter(s)
			server := router.NewServer(r, c.Int(flagPort))

			g, ctx := errgroup.WithContext(ctx)

			g.Go(func() error {
				defer logrus.Info("server exited")
				return server.ListenAndServe()
			})

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			g.Go(func() error {
				defer logrus.Info("signal handler finished")
				select {
				case <-ctx.Done():
					if err := server.Shutdown(ctx); err != nil {
						return err
					}

					return ctx.Err()
				case <-sigChan:
					cancel()
				}
				return nil
			})

			return g.Wait()
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Panic("unable to run app")
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
