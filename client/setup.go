package client

import (
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	pocket "github.com/omgitsotis/pocket-stats/pocket"
)

func NewClient() *Client {
	p := pocket.NewPocket("74935-9d486f66d2999047b61328f3")
	return &Client{Pocket:p}
}

func ServeAPI() error {
	c := NewClient()
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(c.Healthcheck)
	r.Methods("GET").Path("/auth").HandlerFunc(c.Retrieve)
	r.Methods("GET").Path("/auth/recieved").HandlerFunc(c.Authorise)
	r.Methods("GET").Path("/data").HandlerFunc(c.GetData)
	fmt.Println("Created router")
	return http.ListenAndServe(":8080", r)
}
