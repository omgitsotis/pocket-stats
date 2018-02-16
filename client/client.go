package client

import (
	"fmt"
	"net/http"
	"time"

	pocket "github.com/omgitsotis/pocket-stats/pocket"
)

type Client struct {
	Pocket      *pocket.Pocket
	Code        string
	AccessToken string
}

func (c *Client) Retrieve(w http.ResponseWriter, r *http.Request) {
	code, err := c.Pocket.GetAuth("http://localhost:8080/auth/recieved")
	if err != nil {
		WriteErrorResponse(w, err.Error())
		return
	}

	c.Code = code

	u := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		code,
		"http://localhost:8080/auth/recieved",
	)
	fmt.Fprintln(w, u)
}

func (c *Client) Authorise(w http.ResponseWriter, r *http.Request) {
	user, err := c.Pocket.ReceieveAuth(c.Code)
	if err != nil {
		WriteErrorResponse(w, err.Error())
		return
	}

	fmt.Println("user access token", user.AccessToken)

	c.AccessToken = user.AccessToken
	http.SetCookie(w, &http.Cookie{
		Name:  "pocket-token",
		Value: user.AccessToken,
		Path:  "/",
	})

	fmt.Fprintf(w, "%v", user)
}

func (c *Client) GetData(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pocket-token")
	if err != nil {
		WriteErrorResponse(w, err.Error())
		return
	}

	fmt.Println("cookie value", cookie.String())
	fmt.Println("access token", c.AccessToken)

	since := time.Now().AddDate(0, 0, -7).Unix()
	data, err := c.Pocket.GetData(c.AccessToken, since)
	if err != nil {
		WriteErrorResponse(w, err.Error())
		return
	}

	fmt.Fprintf(w, "%v", data)
}

func (c *Client) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pocket-app healthy")
	return
}
