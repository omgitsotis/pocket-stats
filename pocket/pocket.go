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

func (p *Pocket) InitDB(ip model.InitParams) (*model.DataList, error) {
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
			if v.Status == model.Archived {
				wc, err := strconv.Atoi(v.WordCount)
				if err != nil {
					fmt.Printf("Error getting word count %s\n", err.Error())
					continue
				}

				r := model.Row{
					ID:        v.ItemID,
					WordCount: wc,
					DateRead:  t.Unix(),
					Status:    model.Archived,
					UserID:    ip.ID,
				}

				p.dao.AddArticle(r)
				continue
			}

			if v.Status == model.Added {
				wc, err := strconv.Atoi(v.WordCount)
				if err != nil {
					fmt.Printf("Error getting word count %s\n", err.Error())
					continue
				}

				r := model.Row{
					ID:        v.ItemID,
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

	// param := DataParam{
	// 	ConsumerKey: p.ConsumerID,
	// 	AccessToken: token,
	// 	Since:       until,
	// 	State:       "all",
	// 	Sort:        "oldest",
	// }
	//
	var dl model.DataList
	// if err := p.Call("/get", param, &dl); err != nil {
	// 	return nil, err
	// }
	//
	return &dl, nil
}

func (p *Pocket) GetStatsForDates(param model.StatsParams) (*model.Stats, error) {
	added, err := p.dao.GetCountForDates(param.Start, param.End, "date_added")
	if err != nil {
		return nil, err
	}

	read, err := p.dao.GetCountForDates(param.Start, param.End, "date_read")
	if err != nil {
		return nil, err
	}

	stats := model.Stats{added, read}
	return &stats, nil
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
