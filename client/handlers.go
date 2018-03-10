package client

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/omgitsotis/pocket-stats/pocket/model"
)

func sendAuth(client *Client, data interface{}) {
	if client.Code != "" {
		client.send <- Message{"subscribe auth", "authorised"}
		return
	}

	fmt.Println(client.Code)

	code, err := client.Pocket.GetAuth("http://localhost:4000/auth/recieved")
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	client.Code = code

	u := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		code,
		"http://localhost:4000/auth/recieved",
	)

	type Link struct {
		URL string `json:"url"`
	}

	link := Link{u}
	client.send <- Message{"send auth", link}
}

func initDB(client *Client, data interface{}) {
	var params model.InitParams
	err := mapstructure.Decode(data, &params)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	_, err = client.Pocket.InitDB(params)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	client.send <- Message{"data get", "Complete"}
}

func saveToken(client *Client, data interface{}) {
	type AccessToken struct {
		Token string `json:"token"`
	}

	var token AccessToken
	err := mapstructure.Decode(data, &token)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	fmt.Println(data)
	fmt.Println(token)
	client.AccessToken = token.Token

	client.send <- Message{"subscribe auth", client.AccessToken}
}

func getStatistics(client *Client, data interface{}) {
	var p model.StatsParams

	err := mapstructure.Decode(data, &p)
	if err != nil {
		log.Printf("Error decoding params: %s", err.Error())
		client.send <- Message{"error", err.Error()}
		return
	}

	log.Printf("Get stats from %d to %d", p.Start, p.End)

	stats, err := client.Pocket.GetStatsForDates(p)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	client.send <- Message{"data get", stats}
}
