package client

import (
	"fmt"
	"net/http"
	"net/url"
)

type Client struct{}

func (c *Client) Retrieve(w http.ResponseWriter, r *http.Request) {
	// client := &http.Client{}
	r, err := http.NewRequest("POST", "https://getpocket.com/v3/oauth/request", nil)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}

	data := url.Values{}
	data.Set("consumer_key", "74935-9d486f66d2999047b61328f3")
	data.Set("redirect_uri", "localhost:8080/oauth")

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

}
