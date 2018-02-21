package client

import (
	"fmt"
)

func sendAuth(client *Client, data interface{}) {
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
