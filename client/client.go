package client

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	pocket "github.com/omgitsotis/pocket-stats/pocket"
)

type Client struct {
	Pocket *pocket.Pocket
}

func (c *Client) Retrieve(w http.ResponseWriter, r *http.Request) {
	code, err := c.Pocket.GetAuth("http://localhost:8080/auth/recieved")
	if err != nil {
		WriteErrorResponse(w, err.Error())
		return
	}

	u := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		code,
		"http://localhost:8080/auth/recieved",
	)
	fmt.Fprintln(w, u)
}

func (c *Client) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pocket-app healthy")
	return
}

func NewClient() *Client {
	p := pocket.NewPocket("74935-9d486f66d2999047b61328f3")
	return &Client{p}
}

func ServeAPI() error {
	c := NewClient()
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(c.Healthcheck)
	r.Methods("GET").Path("/auth").HandlerFunc(c.Retrieve)
	return http.ListenAndServe(":8080", r)
}
