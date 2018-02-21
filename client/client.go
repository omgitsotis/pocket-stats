package client

import (
	"github.com/gorilla/websocket"
	"github.com/omgitsotis/pocket-stats/pocket"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type FindHandler func(string) (Handler, bool)

type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindHandler
	stopChannels map[int]chan bool
	Pocket       *pocket.Pocket
	Code         string
}

func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

func (c *Client) Read() {
	var message Message
	for {
		if err := c.socket.ReadJSON(&message); err != nil {
			break
		}
		if fn, ok := c.findHandler(message.Name); ok {
			fn(c, message.Data)
		}
	}

	c.socket.Close()
}

func (c *Client) Write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}

	c.socket.Close()
}

func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.send)
}

func NewClient(conn *websocket.Conn, fn FindHandler) *Client {
	return &Client{
		send:         make(chan Message),
		socket:       conn,
		findHandler:  fn,
		stopChannels: make(map[int]chan bool),
		Pocket:       pocket.NewPocket("74935-9d486f66d2999047b61328f3"),
	}
}
