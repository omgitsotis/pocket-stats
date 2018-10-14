package client

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/omgitsotis/pocket-stats/server/pocket"
)

type Handler func(*Client, interface{})

type Router struct {
	rules  map[string]Handler
	pocket *pocket.Pocket
	client *Client
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewRouter(p *pocket.Pocket) *Router {
	return &Router{
		rules:  make(map[string]Handler),
		pocket: p,
	}
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	e.client = NewClient(socket, e.FindHandler, e.pocket)
	defer e.client.Close()
	go e.client.Write()
	e.client.Read()
}

func (r *Router) Handle(name string, fn Handler) {
	r.rules[name] = fn
}

func (r *Router) FindHandler(name string) (Handler, bool) {
	handler, ok := r.rules[name]
	return handler, ok
}

func (e *Router) RecievedAuth(w http.ResponseWriter, r *http.Request) {
	if e.client == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Client not found")
		return
	}

	user, err := e.client.Pocket.ReceieveAuth(e.client.Code)
	if err != nil {
		WriteErrorResponse(w, err.Error())
		return
	}

	e.client.AccessToken = user.AccessToken
	e.client.send <- Message{"auth user", user}

	http.ServeFile(w, r, "./response.html")
}
