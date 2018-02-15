package pocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Pocket struct {
	ConsumerID string
	Client     *http.Client
}

func (p *Pocket) GetAuth(uri string) (string, error) {
	params := Request{p.ConsumerID, uri}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("Error marshalling params: %s", err.Error())
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://getpocket.com/v3/oauth/request",
		bytes.NewBuffer(b),
	)

	if err != nil {
		log.Printf("Error creating request: %s", err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	res, err := p.Client.Do(req)
	if err != nil {
		log.Printf("error performing request: %s", err.Error())
		return "", err
	}

	defer res.Body.Close()

	var rt RequestToken
	err = json.NewDecoder(res.Body).Decode(&rt)
	if err != nil {
		log.Printf("Error decoding body: %s", err.Error())
		return "", err
	}

	log.Printf("repsone code returned [%s]", rt.Code)
	return rt.Code, nil
}

func NewPocket(id string) *Pocket {
	return &Pocket{id, &http.Client{}}
}

func (p *Pocket) SetClient(c http.Client) {
	p.Client = &c
}
