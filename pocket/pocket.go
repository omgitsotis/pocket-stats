package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/omgitsotis/pocket-stats/pocket/dao"
	"github.com/omgitsotis/pocket-stats/pocket/model"
)

var logger = log.New(os.Stdout, "Pocket:", log.Ldate|log.Ltime|log.Lshortfile)

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

	logger.Printf("repsone code returned [%s]", rt.Code)
	return rt.Code, nil
}

func (p *Pocket) ReceieveAuth(key string) (*model.User, error) {
	a := model.Authorise{p.ConsumerID, key}
	var user model.User

	if err := p.Call("/oauth/authorize", a, &user); err != nil {
		return nil, err
	}

	// id, err := p.dao.AddUser(user.Username)
	// if err != nil {
	// 	return nil, err
	// }

	user.ID = 1

	date, err := p.dao.GetLastAdded()
	if err != nil {
		return nil, err
	}

	user.LastUpdated = date

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
				logger.Printf("Error converting ID: %s", err.Error())
				continue
			}

			wc, err := strconv.Atoi(v.WordCount)
			if err != nil {
				fmt.Printf("Error getting word count %s\n", err.Error())
				continue
			}

			if v.Status == model.Archived {
				r := model.Article{
					ID:        int64(id),
					WordCount: int64(wc),
					DateRead:  t.Unix(),
					Status:    model.Archived,
					UserID:    ip.ID,
				}

				p.dao.AddArticle(r)
				continue
			}

			if v.Status == model.Added {
				r := model.Article{
					ID:        int64(id),
					WordCount: int64(wc),
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
	logger.Printf("Current date: %s", midnight.Format("02/01/2006"))
	logger.Printf("Until date: %s", until.Format("02/01/2006"))
	logger.Printf("Days to go back to: %d", days)

	seen := make(map[string]bool)

	for i := 0; i < days; i++ {
		t := midnight.AddDate(0, 0, i*-1)

		logger.Printf("Geting info from date %s", t.Format("02/01/2006"))

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

			logger.Printf("Got article %s\n", v.ItemID)

			id, err := strconv.Atoi(v.ItemID)
			if err != nil {
				logger.Printf("Error converting ID: %s", err.Error())
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
					r := model.Article{
						ID:        int64(id),
						WordCount: int64(wc),
						DateRead:  t.Unix(),
						Status:    model.Archived,
						UserID:    ip.ID,
					}

					p.dao.AddArticle(r)
					continue
				}

				if v.Status == model.Added {
					r := model.Article{
						ID:        int64(id),
						WordCount: int64(wc),
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
					logger.Printf("Article %d has same status [%s], skip", id, v.Status)
				}
			}
		}
	}

	logger.Printf("Updated to %d", midnight.Unix())
	return midnight.Unix(), nil
}

func (p *Pocket) GetStatsForDates(sp model.StatsParams) (*model.Stats, error) {
	articles, err := p.dao.GetArticles(sp.Start, sp.End)
	if err != nil {
		return nil, err
	}

	stats := p.createStats(sp, articles)
	return stats, nil
}

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

func (p *Pocket) GetUser() (*model.User, error) {
	date, err := p.dao.GetLastAdded()
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username:    "omgitsotis",
		ID:          1,
		LastUpdated: date,
	}

	return &user, nil
}

func (p *Pocket) Call(uri string, body, t interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		logger.Printf("Error marshalling params: %s", err.Error())
		return err
	}

	uri = fmt.Sprintf("https://getpocket.com/v3%s", uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		logger.Printf("Error creating request: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		logger.Printf("error performing request: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logger.Printf("Status %s", res.Status)
		return errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		logger.Printf("Error decoding body: %s", err.Error())
		return err
	}

	return nil
}

func (p *Pocket) SetClient(c http.Client) {
	p.Client = &c
}

func (p *Pocket) createStats(sp model.StatsParams, arts []model.Article) *model.Stats {
	stats := make(map[int64]*model.Stat)
	var wAdded, wRead, aAdded, aRead int64
	for _, a := range arts {
		if a.Status == model.Archived {
			// logger.Printf("Read article: words %d | date %d", a.WordCount, a.DateRead)
			aRead++
			wRead += a.WordCount

			s, ok := stats[a.DateRead]
			if ok {
				s.ArticleRead++
				s.WordRead += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleRead: 1,
					WordRead:    a.WordCount,
				}

				stats[a.DateRead] = &newStat
			}

			if a.DateAdded >= sp.Start && a.DateAdded <= sp.End {
				aAdded++
				wAdded += a.WordCount

				s, ok = stats[a.DateAdded]
				if ok {
					s.ArticleAdded++
					s.WordAdded += a.WordCount
				} else {
					newStat := model.Stat{
						ArticleAdded: 1,
						WordAdded:    a.WordCount,
					}

					stats[a.DateAdded] = &newStat
				}
			}

		} else {
			// logger.Printf("Added article: words %d | date %d", a.WordCount, a.DateRead)
			aAdded++
			wAdded += a.WordCount

			s, ok := stats[a.DateAdded]
			if ok {
				s.ArticleAdded++
				s.WordAdded += a.WordCount
			} else {
				newStat := model.Stat{
					ArticleAdded: 1,
					WordAdded:    a.WordCount,
				}

				stats[a.DateAdded] = &newStat
			}
		}
	}

	totals := model.TotalStats{
		ArticlesAdded: aAdded,
		ArticlesRead:  aRead,
		WordsAdded:    wAdded,
		WordsRead:     wRead,
	}

	return &model.Stats{
		Start:  sp.Start,
		End:    sp.End,
		Value:  stats,
		Totals: totals,
	}
}

func NewPocket(id string, c *http.Client, d dao.DAO) *Pocket {
	return &Pocket{id, c, d}
}
