package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pkg/database"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init(l *logrus.Logger) {
	log = l
}

const (
	TagRead  = "_read"
	TagSport = "_sport"
	TagNews  = "_news"
)

type Store interface {
	BeginBatch()
	CommitBatch(ctx context.Context) error
	SaveArticle(a *database.Article)
	DeleteArticle(id string)

	GetArticle(ctx context.Context, id string) (*database.Article, error)
	GetArticlesByDate(ctx context.Context, start, end int64) ([]*database.Article, error)
	GetArticlesByDateAndTag(ctx context.Context, start, end int64, tag string) ([]*database.Article, error)

	GetLastUpdateDate(ctx context.Context) (int, error)
	SaveLastUpdateDate(ctx context.Context, date int) error
}

type Server struct {
	ctx          context.Context
	pocketClient *pocket.Client
	db           Store
	authURL      string
	requestToken string
	username     string
	password     string
}

func New(ctx context.Context, pc *pocket.Client, url string, db Store, user, pass string) *Server {
	return &Server{
		ctx:          ctx,
		pocketClient: pc,
		authURL:      url,
		db:           db,
		username:     user,
		password:     pass,
	}
}

func (s *Server) Healthcheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, healthcheckResp{Status: "running"})
}

// GetAuthLink gets the pocket Authorisation link used to login
func (s *Server) GetAuthLink(w http.ResponseWriter, r *http.Request) {
	code, err := s.pocketClient.GetAuth(s.authURL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting pocket auth code: %w", err))
		return
	}

	s.requestToken = code

	u := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		code,
		s.authURL,
	)

	link := Link{URL: u}
	respondWithJSON(w, http.StatusOK, link)

}

func (s *Server) ReceiveToken(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received request for ReceiveToken")
	_, err := s.pocketClient.ReceieveAuth(s.requestToken)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error getting access token for user: %w", err))
		return
	}

	file, err := ioutil.ReadFile("logged_in.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error reading logged_in.html: %w", err))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
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

	date, err := s.db.GetLastUpdateDate(s.ctx)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error getting last update date: %w", err))
		return
	}

	logrus.Debugf("updating database from [%d]", date)

	response, err := s.pocketClient.GetArticles(date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error getting articles: %w", err))
		return
	}

	articleList := response.GetArticleList()
	if articleList == nil {
		log.Info("No new entries found")
		s.setUpdateDate(w)
		return
	}

	s.db.BeginBatch()
	for _, article := range articleList {
		// If the article has a status of deleted, then delete it in the DB
		if article.Status == pocket.ItemStatusDeleted {
			s.db.DeleteArticle(fmt.Sprintf("%d", article.ItemID))
			continue
		}

		// Otherwise do an upsert
		s.db.SaveArticle(convertArticle(article))
	}

	if err = s.db.CommitBatch(s.ctx); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error updating database: %w", err))
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

	date, err := s.db.GetLastUpdateDate(s.ctx)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error getting last update date: %w", err))
		return
	}

	logrus.Infof("Updating DB from [%d]", date)

	response, err := s.pocketClient.GetArticles(date)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error getting articles: %w", err))
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (s *Server) GetStats(w http.ResponseWriter, r *http.Request) {
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	start, end, err := createDBTime(startParam, endParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	log.Infof("Getting articles between %d - %d", start, end)
	articles, err := s.db.GetArticlesByDate(s.ctx, start, end)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get articles from the DB: %w", err))
		return
	}

	log.Infof("Creating stats for dates %d - %d", start, end)
	totals, err := createTotalStats(start, end, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting total stats from articles: %w", err))
		return
	}

	itemised, err := createItemisedStats(start, end, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting itemised stats from articles: %w", err))
		return
	}

	stats := Stats{
		Totals:   totals,
		Itemised: itemised,
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (s *Server) GetTotalStats(w http.ResponseWriter, r *http.Request) {
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	start, end, err := createDBTime(startParam, endParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	log.Infof("Getting articles between %d - %d", start, end)
	articles, err := s.db.GetArticlesByDate(s.ctx, start, end)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get articles from the DB: %w", err))
		return
	}

	log.Infof("Creating stats for dates %d - %d", start, end)
	totals, err := createTotalStats(start, end, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting total stats from articles: %w", err))
		return
	}

	pStart, pEnd := getPreviousDate(start, end)
	log.Infof("Getting previous articles between %d - %d", start, end)
	articles, err = s.db.GetArticlesByDate(s.ctx, pStart, pEnd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get previous articles from the DB: %w", err))
		return
	}

	log.Infof("Creating previous stats for dates %d - %d", pStart, pEnd)
	previousTotals, err := createTotalStats(pStart, pEnd, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting total stats from articles: %w", err))
		return
	}

	stats := Stats{
		Totals: totals,
		PreviousStats: &PreviousStats{
			Totals: previousTotals,
		},
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (s *Server) GetTagStats(w http.ResponseWriter, r *http.Request) {
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")
	tagParam := r.URL.Query().Get("tags")

	start, end, err := createDBTime(startParam, endParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	var articles []*database.Article

	if tagParam != "" {
		log.Infof("Getting articles between [%d] - [%d] for tag [%s]", start, end, tagParam)
		articles, err = s.db.GetArticlesByDateAndTag(s.ctx, start, end, tagParam)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get articles from the DB: %w", err))
			return
		}
	} else {
		log.Infof("Getting articles between [%d] - [%d]", start, end)
		articles, err = s.db.GetArticlesByDate(s.ctx, start, end)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get articles from the DB: %w", err))
			return
		}
	}

	log.Infof("Creating stats for dates [%d] - [%d] tag [%s]", start, end, tagParam)
	tags, err := createTagStats(start, end, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create tag stats: %w", err))
		return
	}

	pStart, pEnd := getPreviousDate(start, end)
	if tagParam != "" {
		log.Infof("Getting previous articles between [%d] - [%d] for tag [%s]", pStart, pEnd, tagParam)
		articles, err = s.db.GetArticlesByDateAndTag(s.ctx, pStart, pEnd, tagParam)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get previous articles from the DB: %w", err))
			return
		}
	} else {
		log.Infof("Getting previous articles between [%d] - [%d]", pStart, pEnd)
		articles, err = s.db.GetArticlesByDate(s.ctx, pStart, end)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get previous articles from the DB: %w", err))
			return
		}
	}

	log.Infof("Creating previous stats for dates [%d] - [%d] tag [%s]", pStart, pEnd, tagParam)
	previousTags, err := createTagStats(pStart, pEnd, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("error getting total stats from articles: %w", err))
		return
	}

	stats := Stats{
		Tags: tags,
		PreviousStats: &PreviousStats{
			Tags: previousTags,
		},
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (s *Server) GetItemisedStats(w http.ResponseWriter, r *http.Request) {
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	start, end, err := createDBTime(startParam, endParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	log.Infof("Getting articles between [%d] - [%d]", start, end)
	articles, err := s.db.GetArticlesByDate(s.ctx, start, end)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to get articles from the DB: %w", err))
		return
	}

	log.Infof("Creating itemised stats for dates [%d] - [%d]", start, end)
	itemised, err := createItemisedStats(start, end, articles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to create tag stats: %w", err))
		return
	}

	stats := Stats{
		Itemised: itemised,
	}

	respondWithJSON(w, http.StatusOK, stats)
}

func (s *Server) AuthMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			log.Warn("No authorisation provided")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if s.username != user && s.password != pass {
			log.Warn("Incorrect auth info")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fn(w, r)
	}
}

func (s *Server) setUpdateDate(w http.ResponseWriter) {
	updateTime := time.Now().Unix()
	if err := s.db.SaveLastUpdateDate(s.ctx, int(updateTime)); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("error saving last update date: %w", err))
		return
	}

	logrus.Infof("DB updated to [%d]", updateTime)

	respondWithJSON(w, http.StatusOK, updateResponse{Date: updateTime})
}

// ConvertArticle converts the pocket article to the database article
func convertArticle(pa pocket.Article) *database.Article {
	var (
		added, read int
		err         error
	)

	// Convert the added time to a number
	added, err = strconv.Atoi(pa.TimeAdded)
	if err != nil {
		logrus.Errorf(
			"pocket article [%d] has a bad added time [%s]",
			pa.ItemID,
			pa.TimeAdded,
		)
	}

	// Check we have a read time, and then convert that to a number
	if pa.TimeRead != "" {
		read, err = strconv.Atoi(pa.TimeRead)
		if err != nil {
			logrus.Errorf(
				"pocket article [%d] has a bad read time [%s]",
				pa.ItemID,
				pa.TimeRead,
			)
		}
	}

	// Remove the time part of the read and added times
	sAdded := stripTime(added)
	sRead := stripTime(read)

	var tag string
	for key := range pa.Tags {
		if key == TagRead || key == TagSport || key == TagNews {
			continue
		}

		tag = key
		break
	}

	return &database.Article{
		ID:        fmt.Sprintf("%d", pa.ItemID),
		URL:       pa.ResolvedURL,
		Title:     pa.ResolvedTitle,
		WordCount: int64(pa.WordCount),
		DateAdded: sAdded,
		DateRead:  sRead,
		Tag:       tag,
	}
}

// StripTime converts the time part of a time stamp into 00:00:00.
func stripTime(fullTime int) int64 {
	if fullTime == 0 {
		return 0
	}

	t := time.Unix(int64(fullTime), 0)
	stripped := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	return stripped.Unix()
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	respondWithJSON(w, code, APIError{Code: code, Message: err.Error()})
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
