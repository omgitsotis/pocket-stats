package client

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
)

type InitParams struct {
	ID    int    `json:id`
	Token string `json:token`
	Date  int64  `json:"date"`
}

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
	var params InitParams
	err := mapstructure.Decode(data, &params)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	since = time.Unix(params.Date, 0)
	_, err := client.Pocket.InitDB(client.AccessToken, since)
	if err != nil {
		client.send <- Message{"error", err.Error()}
		return
	}

	// client.send <- Message{"data get", data}
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
