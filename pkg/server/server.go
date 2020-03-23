package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pkg/database"
	"github.com/omgitsotis/pocket-stats/pkg/pocket"
	"github.com/pkg/errors"
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
	username     string
	password     string
}

func New(pc *pocket.Client, url string, db database.DBCLient, user, pass string) *Server {
	return &Server{
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

	link := Link{URL: u}
	respondWithJSON(w, http.StatusOK, link)

}

func (s *Server) ReceiveToken(w http.ResponseWriter, r *http.Request) {
	log.Debug("Received request for ReceiveToken")
	_, err := s.pocketClient.ReceieveAuth(s.requestToken)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"error getting access token for user",
			err,
		)
	}

	file, err := ioutil.ReadFile("logged_in.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error reading logged_in.html", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
	// respondWithJSON(w, http.StatusOK, user)
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

	if err = s.db.UpsertArticles(articleList); err != nil {
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

func (s *Server) GetStats(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	if start == "" {
		respondWithError(w, http.StatusBadRequest, "No start date provided", nil)
		return
	}

	if end == "" {
		respondWithError(w, http.StatusBadRequest, "No end date provided", nil)
		return
	}

	startInt, err := strconv.Atoi(start)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Could not convert start date [%s]", start),
			err,
		)
		return
	}

	endInt, err := strconv.Atoi(end)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("Could not convert end date [%s]", end),
			err,
		)
		return
	}

	dbStartTime := database.StripTime(startInt)
	dbEndTime := database.StripTime(endInt)

	log.Infof("Getting articles between %d - %d", dbStartTime, dbEndTime)
	articles, err := s.db.GetArticlesByDate(dbStartTime, dbEndTime)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"unable to get articles from the DB",
			err,
		)
		return
	}

	log.Infof("Creating stats for dates %d - %d", dbStartTime, dbEndTime)
	stats, err := createStats(dbStartTime, dbEndTime, articles)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Error converting articles to stats",
			err,
		)
		return
	}

	log.Infof("Created stats for dates %d - %d", dbStartTime, dbEndTime)
	respondWithJSON(w, http.StatusOK, stats)
}

func (s *Server) setUpdateDate(w http.ResponseWriter) {
	updateTime := time.Now().Unix()
	if err := s.db.SaveUpdateDate(updateTime); err != nil {
		respondWithError(w, http.StatusBadRequest, "error saving last update date", err)
	}

	logrus.Infof("DB updated to [%d]", updateTime)

	respondWithJSON(w, http.StatusOK, updateResponse{Date: updateTime})
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

func createStats(start, end int64, articles []database.Article) (*Stats, error) {
	itemised := make(ItemisedStats)
	tags := make(TagStats)
	st := StatTotals{}

	// Populate the itemised map. We want all the dates in the range, including
	// the days with no updates
	t := time.Unix(start, 0)
	endTime := time.Unix(end, 0)

	for {
		itemised[t.Unix()] = &StatTotals{}
		t = t.AddDate(0, 0, 1)
		log.Debugf("%d greater than %d", t.Unix(), endTime.Unix())
		if t.Unix() > endTime.Unix() {
			break
		}
	}

	for _, a := range articles {
		log.Debugf("Checking article [%d]", a.ID)
		// Check to see if the article was added in the date range
		if isInRange(start, end, a.DateAdded) {
			log.Debugf(
				"Article [%d]: Start date [%d] < Article add date [%d] < End date [%d]",
				a.ID, start, a.DateAdded, end,
			)

			dayAddedTotal, ok := itemised[a.DateAdded]
			if !ok {
				return nil, errors.Errorf(
					"What the fuck, date [%d] not created",
					a.DateAdded,
				)
			}

			timeReading := convertWordsToTime(a.WordCount)
			// Update itemised values
			dayAddedTotal.ArticlesAdded++
			dayAddedTotal.WordsAdded += a.WordCount
			dayAddedTotal.TimeAdded += timeReading

			// Update total values
			st.ArticlesAdded++
			st.WordsAdded += a.WordCount
			st.TimeAdded += timeReading
		}

		// Check to see if the article is read
		if a.DateRead != 0 && isInRange(start, end, a.DateRead) {
			log.Debugf("Article [%d] read [%d]", a.ID, a.DateRead)
			dayReadTotal, ok := itemised[a.DateRead]
			if !ok {
				return nil, errors.Errorf(
					"What the fuck, date [%d] not created",
					a.DateRead,
				)
			}

			timeReading := convertWordsToTime(a.WordCount)
			// Update itemised values
			dayReadTotal.ArticlesRead++
			dayReadTotal.WordsRead += a.WordCount
			dayReadTotal.TimeRead += timeReading

			// Update total values
			st.ArticlesRead++
			st.WordsRead += a.WordCount
			st.TimeRead += timeReading

			// Update the tag values
			if _, ok := tags[a.Tag]; !ok {
				tags[a.Tag] = &StatTotals{
					ArticlesRead: 1,
					WordsRead:    a.WordCount,
					TimeRead:     timeReading,
				}
			} else {
				tags[a.Tag].ArticlesRead++
				tags[a.Tag].WordsRead += a.WordCount
				tags[a.Tag].TimeRead += timeReading
			}
		}
	}

	return &Stats{
		Totals:   st,
		Itemised: itemised,
		Tags:     tags,
	}, nil
}

func convertWordsToTime(words int64) int64 {
	timeReading := float64(words / WordsPerMinute)
	rounded := math.Round(timeReading)
	return int64(rounded)
}

// isInRange checks to see if the article was added within the date range
func isInRange(start, end, added int64) bool {
	return start <= added && added <= end
}

func respondWithError(w http.ResponseWriter, code int, message string, err error) {
	if err != nil {
		log.WithError(err).Error(message)
	}

	respondWithJSON(w, code, APIError{Code: code, Message: message})
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
