package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/omgitsotis/pocket-stats/pkg/server"
	"github.com/sirupsen/logrus"
)

const route = "/api/pocket/"

var log *logrus.Logger

func Init(l *logrus.Logger) {
	log = l
}

func CreateRouter(s *server.Server) *mux.Router {
	router := mux.NewRouter()

	router.NewRoute().
		Path(route).
		Methods(http.MethodGet).
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hit app"))
			},
		)

	router.NewRoute().
		Path(route + "auth").
		Methods(http.MethodGet).
		HandlerFunc(s.GetAuthLink())

	router.NewRoute().
		Path(route + "auth/received").
		Methods(http.MethodGet).
		HandlerFunc(s.ReceiveToken())

	return router
}

func NewServer(router *mux.Router, port int) *http.Server {
	allowedCorsMethods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPost,
		http.MethodOptions,
	})
	allowedCorsOrigins := handlers.AllowedOrigins([]string{"*"})

	return &http.Server{
		Handler:      handlers.CORS(allowedCorsMethods, allowedCorsOrigins)(router),
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Duration(15) * time.Second,
		ReadTimeout:  time.Duration(15) * time.Second,
	}
}

func StartServer(server *http.Server) {
	log.Info("Pocket stats API starting with address " + server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("Fatal error while running HTTP server")
	}
}
