package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/omgitsotis/pocket-stats/pkg/database"
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
	db           database.DBCLient
	authURL      string
	requestToken string
}

func New(pc *pocket.Client, url string, db database.DBCLient) *Server {
	return &Server{
		pocketClient: pc,
		authURL:      url,
		db:           db,
	}
}

// GetAuthLink gets the pocket Authorisation link used to login
func (s *Server) GetAuthLink(w http.ResponseWriter, r *http.Request) {
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

	link := model.Link{URL: u}
	respondWithJSON(w, http.StatusOK, link)

}

func (s *Server) ReceiveToken(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {

	log.Debug("Received request for CheckAuthStatus")

	if s.pocketClient.IsAuthed() {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	return
}

func (s *Server) GetArticles(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received request for GetArticles")
	if !s.pocketClient.IsAuthed() {
		respondWithError(w, http.StatusForbidden, "User not authorised", nil)
		return
	}

	// Run the database update in the background
	go s.getArticles()

	respondWithJSON(w, http.StatusOK, &pocket.RetrieveResult{Complete: 1})
}

func (s *Server) getArticles() {
	complete := false
	index := 0
	for !complete {
		log.Infof("Getting articles [%d]-[%d]", (100 * index), (100 * (index + 1)))
		resp, err := s.pocketClient.GetArticles(100 * index)
		if err != nil {
			log.WithError(err).Error("Error getting articles")
			return
		}

		articles := resp.GetArticleList()

		if err = s.db.SaveArticles(articles); err != nil {
			log.WithError(err).Error("Error saving articles")
			return
		}

		if resp.Complete == 1 || index == 2 {
			complete = true
		}

		index++
	}
	log.Info("Finished updating database")
}

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	if err != nil {
		log.WithError(err).Error(message)
	}

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
