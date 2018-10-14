package client

import (
	"net/http"
	"os"

	"github.com/omgitsotis/pocket-stats/server/pocket"
	"github.com/omgitsotis/pocket-stats/server/pocket/dao/sqlite"
	logging "github.com/op/go-logging"
)

var clientLog = logging.MustGetLogger("client")

// Serve API handles the creation of the logging, database and router and then
// runs the server
func ServeAPI() error {
	setupLogging()

	sqlite, err := sqlite.NewSQLiteDAO("./database/pocket.db")
	if err != nil {
		return err
	}

	p := pocket.NewPocket(
		"74935-9d486f66d2999047b61328f3",
		&http.Client{},
		sqlite,
	)

	r := createRouter(p)

	http.Handle("/", r)
	// Create the end point for pocket to return a response after authenticating
	http.HandleFunc("/auth/recieved", r.RecievedAuth)

	clientLog.Info("Serving application")
	return http.ListenAndServe(":4000", nil)

}

// setupLogging handles the creation and format for the logs. Outputs to console
// for now
func setupLogging() {
	var format = logging.MustStringFormatter(
		`%{color}[%{time:Mon 02 Jan 2006 15:04:05.000}] %{level:.5s} %{shortfile} %{color:reset} %{message}`,
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formatter)
}

// createRouter handles the creation of the router and all of the required
// routes
func createRouter(p *pocket.Pocket) *Router {
	r := NewRouter(p)
	r.Handle("send auth", sendAuth)
	r.Handle("data init", initDB)
	r.Handle("auth cached", loadUser)
	r.Handle("data get", getStatistics)
	r.Handle("data update", updateDB)
	r.Handle("data load", loadData)

	return r
}
