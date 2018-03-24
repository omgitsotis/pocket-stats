package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pocket/dao"
	"github.com/omgitsotis/pocket-stats/pocket/model"
)

// Pocket is the client to the pocket api
type Pocket struct {
	ConsumerID string
	Client     *http.Client
	dao        dao.DAO
}

// GetAuth gets the request token from pocket
func (p *Pocket) GetAuth(uri string) (string, error) {
	r := model.Request{p.ConsumerID, uri}
	var rt model.RequestToken
	if err := p.call("/oauth/request", r, &rt); err != nil {
		return "", err
	}

	log.Debugf("repsone code returned [%s]", rt.Code)
	return rt.Code, nil
}

// RecievedAuth gets the access token from pocket, and returns the user from the
// database
func (p *Pocket) ReceieveAuth(key string) (*model.User, error) {
	a := model.Authorise{p.ConsumerID, key}
	var user model.User

	if err := p.call("/oauth/authorize", a, &user); err != nil {
		return nil, err
	}

	user.ID = 1

	date, err := p.dao.GetLastAdded()
	if err != nil {
		return nil, err
	}

	user.LastUpdated = date
	log.Infof("Last added date: [%d]", date)

	return &user, nil
}

// InitDB loads the database with the user's pocket information from a given
// date
func (p *Pocket) InitDB(ip model.InputParams) error {
	ok, err := p.dao.IsUser(ip.ID)
	if err != nil {
		return err
	}

	if !ok {
		log.Error("No user found to init DB")
		return errors.New("No user id found")
	}

	until := time.Unix(ip.Date, 0)

	year, month, day := time.Now().Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	d := midnight.Sub(until)

	days := int(d.Hours() / 24)
	log.Debugf("Current date: %s", midnight.Format("02/01/2006"))
	log.Debugf("Until date: %s", until.Format("02/01/2006"))
	log.Debugf("Days to go back to: %d", days)

	seen := make(map[string]bool)
	var dl model.DataList

	param := model.DataParam{
		ConsumerKey: p.ConsumerID,
		AccessToken: ip.Token,
		State:       "all",
		Sort:        "oldest",
		Type:        "simple",
	}

	for i := 0; i < days; i++ {
		unixTime := midnight.AddDate(0, 0, i*-1).Unix()
		param.Since = unixTime

		if err := p.call("/get", param, &dl); err != nil {
			return err
		}

		for k, v := range dl.Values {
			if seen[k] {
				continue
			}

			seen[k] = true

			if v.Status == model.Deleted {
				continue
			}

			p.addArticle(v, ip.ID, unixTime)
		}
	}

	log.Info("Finished init")
	return nil
}

// UpdateDB updates the database from the last update date to now
func (p *Pocket) UpdateDB(ip model.InputParams) (int64, error) {
	// Get the last updated date
	date, err := p.dao.GetLastAdded()
	if err != nil {
		return 0, err
	}

	until := time.Unix(date, 0)

	year, month, day := time.Now().Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	d := midnight.Sub(until)

	days := int(d.Hours() / 24)
	log.Debugf("Current date: %s", midnight.Format("02/01/2006"))
	log.Debugf("Until date: %s", until.Format("02/01/2006"))
	log.Debugf("Days to go back to: %d", days)

	seen := make(map[string]bool)

	param := model.DataParam{
		ConsumerKey: p.ConsumerID,
		AccessToken: ip.Token,
		State:       "all",
		Sort:        "oldest",
		Type:        "simple",
	}

	for i := 0; i < days+1; i++ {
		t := midnight.AddDate(0, 0, i*-1)
		param.Since = t.Unix()

		log.Debugf("Geting info from date %s", t.Format("02/01/2006"))

		var dl model.DataList
		if err := p.call("/get", param, &dl); err != nil {
			return 0, err
		}

		for k, v := range dl.Values {
			if seen[k] {
				continue
			}

			seen[k] = true

			if v.Status == model.Deleted {
				continue
			}

			id, err := strconv.Atoi(v.ItemID)
			if err != nil {
				log.Errorf("Error converting ID: %s", err.Error())
				continue
			}

			row, err := p.dao.GetArticle(int64(id))
			if err != nil {
				continue
			}

			// Insert
			if row == nil {
				wc, err := strconv.Atoi(v.WordCount)
				if err != nil {
					log.Errorf("Error getting word count %s", err.Error())
					continue
				}

				r := model.Article{
					ID:        int64(id),
					WordCount: int64(wc),
					Status:    v.Status,
					UserID:    ip.ID,
				}

				if v.Status == model.Archived {
					r.DateRead = t.Unix()
					log.Debugf("Adding read article %d", id)
				}

				if v.Status == model.Added {
					r.DateAdded = t.Unix()
					log.Debugf("Adding read article %d", id)
				}

				p.dao.AddArticle(r)
			} else {
				if v.Status != row.Status {
					row.DateRead = t.Unix()
					row.Status = model.Archived
					log.Debugf("Marking article %d as read", id)
					p.dao.UpdateArticle(row)
				}
			}
		}
	}

	log.Infof("Updated to %d", midnight.Unix())
	return midnight.Unix(), nil
}

// GetStatsForDates returns basic stats for all articles between the given
// dates
func (p *Pocket) GetStatsForDates(sp model.StatsParams) (*model.Stats, error) {
	articles, err := p.dao.GetArticles(sp.Start, sp.End)
	if err != nil {
		return nil, err
	}

	stats := createStats(sp, articles)
	return stats, nil
}

// LoadData gets basic stats for articles added in the past week.
func (p *Pocket) LoadData() (*model.Stats, error) {
	date, err := p.dao.GetLastAdded()
	if err != nil {
		return nil, err
	}

	to := time.Unix(date, 0)
	from := to.AddDate(0, 0, -7)

	param := model.StatsParams{
		from.Unix(),
		date,
	}

	stats, err := p.GetStatsForDates(param)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetUser gets the user for a given ID
func (p *Pocket) GetUser() (*model.User, error) {
	date, err := p.dao.GetLastAdded()
	if err != nil {
		return nil, err
	}

	log.Infof("Last added date: [%d]", date)

	user := model.User{
		Username:    "omgitsotis",
		ID:          1,
		LastUpdated: date,
	}

	return &user, nil
}

// addArticle adds an article to the database
func (p *Pocket) addArticle(d model.Data, userID, date int64) {
	id, err := strconv.Atoi(d.ItemID)
	if err != nil {
		log.Warningf("Error converting ID: %s", err.Error())
		return
	}

	wc, err := strconv.Atoi(d.WordCount)
	if err != nil {
		log.Warningf("Error getting word count %s", err.Error())
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
		log.Debugf("Adding read article %d", id)
	}

	if d.Status == model.Added {
		article.DateAdded = date
		log.Debugf("Adding unread article %d", id)
	}

	p.dao.AddArticle(article)
}

// call makes api requests to the Pocket api and marshal the results.
func (p *Pocket) call(uri string, body, t interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		log.Errorf("Error marshalling params: %s", err.Error())
		return err
	}

	uri = fmt.Sprintf("https://getpocket.com/v3%s", uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Errorf("Error creating request: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		log.Errorf("error performing request: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Debugf("Status %s", res.Status)
		return errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		log.Errorf("Error decoding body: %s", err.Error())
		return err
	}

	return nil
}

// SetClient sets the client used to connect to pocket
func (p *Pocket) SetClient(c http.Client) {
	p.Client = &c
}
