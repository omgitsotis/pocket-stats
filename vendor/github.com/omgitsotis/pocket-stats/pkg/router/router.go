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

var log *logrus.Logger

func Init(l *logrus.Logger) {
	log = l
}

func CreateRouter(s *server.Server) *mux.Router {
	router := mux.NewRouter()

	router.NewRoute().
		Path("/").
		Methods(http.MethodGet).
		HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hit app"))
			},
		)

	sub := router.PathPrefix("/api/pocket").Subrouter()
	sub.NewRoute().
		Path("/auth").
		Methods(http.MethodGet).
		HandlerFunc(s.GetAuthLink)

	sub.NewRoute().
		Path("/auth/received").
		Methods(http.MethodGet).
		HandlerFunc(s.ReceiveToken)

	sub.NewRoute().
		Path("/auth/authed").
		Methods(http.MethodGet).
		HandlerFunc(s.CheckAuthStatus)

	sub.NewRoute().
		Path("/update").
		Methods(http.MethodGet).
		HandlerFunc(s.UpdateArticle)

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
