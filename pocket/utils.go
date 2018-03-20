package pocket

import (
	"net/http"
	"os"

	"github.com/omgitsotis/pocket-stats/pocket/dao"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("pocket")

func setupLogging() {
	var format = logging.MustStringFormatter(
		`%{color}[%{time:Mon 02 Jan 2006 15:04:05.000}] %{level:.5s} %{shortfile} %{color:reset} %{message}`,
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formatter)
}

func NewPocket(id string, c *http.Client, d dao.DAO) *Pocket {
	setupLogging()
	return &Pocket{id, c, d}
}
