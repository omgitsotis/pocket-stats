package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func (p *Pocket) GetData(token string, since int64) (*DataList, error) {
	param := DataParam{
		ConsumerKey: p.ConsumerID,
		AccessToken: token,
		Since:       since,
		State:       "all",
	}

	var d DataList
	if err := p.Call("/get", param, &d); err != nil {
		return nil, err
	}

	return &d, nil
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
