package client

import (
    "github.com/gorilla/websocket"
    "net/http"
    "fmt"
    r "github.com/dancannon/gorethink"
)

type Handler func(*Client, interface{})

type Router struct {
    rules map[string]Handler
    session *r.Session
    client *Client
}

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func (r *http.Request) bool {return true},
}

func NewRouter(session *r.Session) *Router {
    return &Router {
        rules : make(map[string]Handler),
        session: session,
    }
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    socket, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, err)
        return
    }

    c := NewClient(socket, e.FindHandler, e.session)
    e.client = c
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
    e.client.send <- Message{"subscribe auth", user.AccessToken}
}
