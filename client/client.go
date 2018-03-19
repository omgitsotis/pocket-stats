package client

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/omgitsotis/pocket-stats/pocket"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Error struct {
	Msg      string `json:"msg"`
	HookName string `json:"hookname"`
}

type FindHandler func(string) (Handler, bool)

type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindHandler
	stopChannels map[int]chan bool
	Pocket       *pocket.Pocket
	Code         string
	AccessToken  string
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
		mType, r, err := c.socket.NextReader()
		if err != nil {
			log.Printf("Error getting reader %s", err.Error())
			break
		}

		log.Printf("Message type %d", mType)
		d := json.NewDecoder(r)
		d.UseNumber()
		if err := d.Decode(&message); err != nil {
			log.Fatal(err)
			break
		}

		// if err := c.socket.ReadJSON(&message); err != nil {
		// 	break
		// }

		if fn, ok := c.findHandler(message.Name); ok {
			log.Printf("Received request for %s\n", message.Name)
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

func (c *Client) SendError(name string, err error) {
	e := Error{
		HookName: name,
		Msg:      err.Error(),
	}

	c.send <- Message{"error", e}
}

func NewClient(conn *websocket.Conn, fn FindHandler, p *pocket.Pocket) *Client {
	return &Client{
		send:         make(chan Message),
		socket:       conn,
		findHandler:  fn,
		stopChannels: make(map[int]chan bool),
		Pocket:       p,
	}
}
