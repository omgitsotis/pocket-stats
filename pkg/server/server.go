package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/omgitsotis/pocket-stats/pkg/model"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(l *logrus.Logger) {
	log = l
}

type Server struct {
	pocketClient *pocket.Client
	authURL      string
	requestToken string
}

func New(pc *pocket.Client, url string) *Server {
	return &Server{
		pocketClient: pc,
		authURL:      url,
	}
}

// GetAuth gets the pocket Authorisation link used to login
func (s *Server) GetAuthLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := s.pocketClient.GetAuth(s.authURL)
		if err != nil {
			respondWithError(
				w,
				http.StatusInternalServerError,
				"error getting pocket auth code",
				err,
			)
			return
		}

		s.requestToken = code

		u := fmt.Sprintf(
			"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
			code,
			s.authURL,
		)

		link := model.Link{u}
		respondWithJSON(w, http.StatusOK, link)
	}
}

func (s *Server) ReceiveToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Received request for ReceiveToken")
		user, err := s.pocketClient.ReceieveAuth(s.requestToken)
		if err != nil {
			respondWithError(
				w,
				http.StatusBadRequest,
				"error getting access token for user",
				err,
			)
		}

		respondWithJSON(w, http.StatusOK, user)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	log.WithError(err).Error(message)
	respondWithJSON(w, code, model.APIError{Code: code, Message: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Info("error marshaling data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
