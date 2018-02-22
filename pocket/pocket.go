package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"
)

type Pocket struct {
	ConsumerID string
	Client     *http.Client
}

func (p *Pocket) GetAuth(uri string) (string, error) {
	r := Request{p.ConsumerID, uri}
	var rt RequestToken
	if err := p.Call("/oauth/request", r, &rt); err != nil {
		return "", err
	}

	log.Printf("repsone code returned [%s]", rt.Code)
	return rt.Code, nil
}

func (p *Pocket) ReceieveAuth(key string) (*User, error) {
	a := Authorise{p.ConsumerID, key}
	var user User

	if err := p.Call("/oauth/authorize", a, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *Pocket) GetData(token string, until time.Time) (*DataList, error) {
	year, month, day := time.Now().Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	d := midnight.Sub(until)

	days := int(d.Hours() / 24)
	fmt.Printf("Current date: %s\n", midnight.Format("02/01/2006"))
	fmt.Printf("Until date: %s\n", until.Format("02/01/2006"))
	fmt.Printf("Days to go back to: %d\n", days)

	seen := make(map[string]bool)
	// data

	for i := 0; i < 2; i++ {
		t := midnight.AddDate(0, 0, i*-1)
		param := DataParam{
			ConsumerKey: p.ConsumerID,
			AccessToken: token,
			Since:       t.Unix(),
			State:       "all",
			Sort:        "oldest",
		}

		var dl DataList
		if err := p.Call("/get", param, &dl); err != nil {
			return nil, err
		}

		count := 0
		s := make([]string, 0)
		for key, _ := range dl.Values {
			fmt.Printf("ID: %s\n", key)
			if seen[key] {
				continue
			}
			s = append(s, key)
			seen[key] = true
			count++
		}

		sort.Strings(s)
		fmt.Println(s)
	}

	// param := DataParam{
	// 	ConsumerKey: p.ConsumerID,
	// 	AccessToken: token,
	// 	Since:       until,
	// 	State:       "all",
	// 	Sort:        "oldest",
	// }
	//
	var dl DataList
	// if err := p.Call("/get", param, &dl); err != nil {
	// 	return nil, err
	// }
	//
	return &dl, nil
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

func NewPocket(id string) *Pocket {
	return &Pocket{id, &http.Client{}}
}

func (p *Pocket) SetClient(c http.Client) {
	p.Client = &c
}
