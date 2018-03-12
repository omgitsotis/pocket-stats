package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pocket/dao"
	"github.com/omgitsotis/pocket-stats/pocket/model"
)

type Pocket struct {
	ConsumerID string
	Client     *http.Client
	dao        dao.DAO
}

func (p *Pocket) GetAuth(uri string) (string, error) {
	r := model.Request{p.ConsumerID, uri}
	var rt model.RequestToken
	if err := p.Call("/oauth/request", r, &rt); err != nil {
		return "", err
	}

	log.Printf("repsone code returned [%s]", rt.Code)
	return rt.Code, nil
}

func (p *Pocket) ReceieveAuth(key string) (*model.User, error) {
	a := model.Authorise{p.ConsumerID, key}
	var user model.User

	if err := p.Call("/oauth/authorize", a, &user); err != nil {
		return nil, err
	}

	id, err := p.dao.AddUser(user.Username)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return &user, nil
}

func (p *Pocket) InitDB(ip model.InputParams) (*model.DataList, error) {
	ok, err := p.dao.IsUser(ip.ID)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("No user id found")
	}

	until := time.Unix(ip.Date, 0)

	year, month, day := time.Now().Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	d := midnight.Sub(until)

	days := int(d.Hours() / 24)
	fmt.Printf("Current date: %s\n", midnight.Format("02/01/2006"))
	fmt.Printf("Until date: %s\n", until.Format("02/01/2006"))
	fmt.Printf("Days to go back to: %d\n", days)

	seen := make(map[string]bool)

	for i := 0; i < days; i++ {
		t := midnight.AddDate(0, 0, i*-1)

		param := model.DataParam{
			ConsumerKey: p.ConsumerID,
			AccessToken: ip.Token,
			Since:       t.Unix(),
			State:       "all",
			Sort:        "oldest",
			Type:        "simple",
		}

		var dl model.DataList
		if err := p.Call("/get", param, &dl); err != nil {
			return nil, err
		}

		for k, v := range dl.Values {
			fmt.Printf("ID: %s\n", k)
			if seen[k] {
				continue
			}

			seen[k] = true

			if v.Status == model.Deleted {
				continue
			}

			fmt.Printf("Got article %#v\n", v)

			id, err := strconv.Atoi(v.ItemID)
			if err != nil {
				log.Printf("Error converting ID: %s", err.Error())
				continue
			}

			wc, err := strconv.Atoi(v.WordCount)
			if err != nil {
				fmt.Printf("Error getting word count %s\n", err.Error())
				continue
			}

			if v.Status == model.Archived {
				r := model.Row{
					ID:        int64(id),
					WordCount: wc,
					DateRead:  t.Unix(),
					Status:    model.Archived,
					UserID:    ip.ID,
				}

				p.dao.AddArticle(r)
				continue
			}

			if v.Status == model.Added {
				r := model.Row{
					ID:        int64(id),
					WordCount: wc,
					DateAdded: t.Unix(),
					Status:    model.Added,
					UserID:    ip.ID,
				}

				p.dao.AddArticle(r)
				continue
			}
		}
	}

	var dl model.DataList
	return &dl, nil
}

func (p *Pocket) GetStatsForDates(param model.StatsParams) (*model.Stats, error) {
	stats, err := p.dao.GetCountForDates(param.Start, param.End)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

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
	log.Printf("Current date: %s", midnight.Format("02/01/2006"))
	log.Printf("Until date: %s", until.Format("02/01/2006"))
	log.Printf("Days to go back to: %d", days)

	seen := make(map[string]bool)

	for i := 0; i < days; i++ {
		t := midnight.AddDate(0, 0, i*-1)

		log.Printf("Geting info from date %s", t.Format("02/01/2006"))

		param := model.DataParam{
			ConsumerKey: p.ConsumerID,
			AccessToken: ip.Token,
			Since:       t.Unix(),
			State:       "all",
			Sort:        "oldest",
			Type:        "simple",
		}

		var dl model.DataList
		if err := p.Call("/get", param, &dl); err != nil {
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

			log.Printf("Got article %s\n", v.ItemID)

			id, err := strconv.Atoi(v.ItemID)
			if err != nil {
				log.Printf("Error converting ID: %s", err.Error())
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
					fmt.Printf("Error getting word count %s\n", err.Error())
					continue
				}

				if v.Status == model.Archived {
					r := model.Row{
						ID:        int64(id),
						WordCount: wc,
						DateRead:  t.Unix(),
						Status:    model.Archived,
						UserID:    ip.ID,
					}

					p.dao.AddArticle(r)
					continue
				}

				if v.Status == model.Added {
					r := model.Row{
						ID:        int64(id),
						WordCount: wc,
						DateAdded: t.Unix(),
						Status:    model.Added,
						UserID:    ip.ID,
					}

					p.dao.AddArticle(r)
					continue
				}
			} else {
				// Update row
				if v.Status != row.Status {
					row.DateRead = t.Unix()
					row.Status = model.Archived
					p.dao.UpdateArticle(row)
				} else {
					log.Printf("Article %d has same status [%s], skip", id, v.Status)
				}
			}
		}
	}

	log.Printf("Updated to %d", midnight.Unix())
	return midnight.Unix(), nil
}

func (p *Pocket) Call(uri string, body, t interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshalling params: %s", err.Error())
		return err
	}

	uri = fmt.Sprintf("https://getpocket.com/v3%s", uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Error creating request: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		log.Printf("error performing request: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Status %s", res.Status)
		return errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		log.Printf("Error decoding body: %s", err.Error())
		return err
	}

	return nil
}

func NewPocket(id string, c *http.Client, d dao.DAO) *Pocket {
	return &Pocket{id, c, d}
}

func (p *Pocket) SetClient(c http.Client) {
	p.Client = &c
}
