package client

import (
    "github.com/gorilla/websocket"
    "net/http"
    "fmt"
)

type Handler func(*Client, interface{})

type Router struct {
    rules map[string]Handler
    client *Client
}

var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func (r *http.Request) bool {return true},
}

func NewRouter() *Router {
    return &Router {
        rules : make(map[string]Handler),
    }
}

func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    socket, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprint(w, err)
        return
    }

    c := NewClient(socket, e.FindHandler)
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

    e.client.send <- Message{"subscribe auth", "authorised"}
}
