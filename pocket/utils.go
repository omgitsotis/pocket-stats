package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/omgitsotis/pocket-stats/server/pocket/dao"
	"github.com/omgitsotis/pocket-stats/server/pocket/model"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("pocket")

const READ_NOW = "read now"
const ATLANTA = "atlanta"
const EXIT_SURVEY = "exit_survey"

func setupLogging() {
	var format = logging.MustStringFormatter(
		`%{color}[%{time:Mon 02 Jan 2006 15:04:05.000}] %{level:.5s} %{shortfile} %{color:reset} %{message}`,
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(formatter)
}

// addArticle adds an article to the database
func (p *Pocket) addArticle(d model.Data, userID, date int64) {
	id, err := strconv.Atoi(d.ItemID)
	if err != nil {
		logger.Warningf("Error converting ID: %s", err.Error())
		return
	}

	wc, err := strconv.Atoi(d.WordCount)
	if err != nil {
		logger.Warningf("Error getting word count %s", err.Error())
		return
	}

	article := model.Article{
		ID:        int64(id),
		WordCount: int64(wc),
		Status:    d.Status,
		UserID:    userID,
	}

	if d.Status == model.Archived {
		article.DateRead = date
		logger.Debugf("Adding read article %d", id)
	}

	if d.Status == model.Added {
		article.DateAdded = date
		logger.Debugf("Adding unread article %d", id)
	}

	p.dao.AddArticle(article)
}

// call makes api requests to the Pocket api and marshal the results.
func (p *Pocket) call(uri string, body, t interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		logger.Errorf("Error marshalling params: %s", err.Error())
		return err
	}

	uri = fmt.Sprintf("https://getpocket.com/v3%s", uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		logger.Errorf("Error creating request: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		logger.Errorf("error performing request: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Debugf("Status %s", res.Status)
		return errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		logger.Errorf("Error decoding body: %s", err.Error())
		return err
	}

	return nil
}

func getTag(tags map[string]model.Tag) string {
	if len(tags) == 0 {
		return ""
	}

	var dbTag string
	for tag, _ := range tags {
		if tag == READ_NOW || tag == EXIT_SURVEY || tag == ATLANTA {
			continue
		}

		dbTag = tag
	}

	logger.Infof("Using tag %s", dbTag)
	return dbTag
}

func NewPocket(id string, c *http.Client, d dao.DAO) *Pocket {
	setupLogging()
	return &Pocket{id, c, d}
}
