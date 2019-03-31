package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/omgitsotis/pocket-stats/pkg/database"
	"github.com/omgitsotis/pocket-stats/pkg/model"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(l *logrus.Logger) {
	log = l
}

type updateResponse struct {
	Date int64 `json:"date_updated"`
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

func (s *Server) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if !s.pocketClient.IsAuthed() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	date, err := s.db.GetLastUpdateDate()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting last update date", err)
		return
	}

	logrus.Infof("Updating DB from [%d]", date)

	response, err := s.pocketClient.GetArticles(date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting articles", err)
		return
	}

	articleList := response.GetArticleList()
	if articleList == nil {
		log.Info("No new entries found")
		s.setUpdateDate(w)
		return
	}

	if err = s.db.UpdateArticles(articleList); err != nil {
		respondWithError(w, http.StatusBadRequest, "error updating articles", err)
		return
	}

	s.setUpdateDate(w)
	return

}

func (s *Server) DebugGetArticle(w http.ResponseWriter, r *http.Request) {
	if !s.pocketClient.IsAuthed() {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	date, err := s.db.GetLastUpdateDate()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting last update date", err)
		return
	}

	logrus.Infof("Updating DB from [%d]", date)

	response, err := s.pocketClient.GetArticles(date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error getting articles", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (s *Server) setUpdateDate(w http.ResponseWriter) {
	updateTime := time.Now().Unix()
	if err := s.db.SaveUpdateDate(updateTime); err != nil {
		respondWithError(w, http.StatusBadRequest, "error saving last update date", err)
	}

	logrus.Infof("DB updated to [%d]", updateTime)

	respondWithJSON(w, http.StatusOK, updateResponse{Date: updateTime})
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
