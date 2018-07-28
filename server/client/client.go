package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/omgitsotis/pocket-stats/server/pocket"
)

// Message is the object that is returned when sending a response via websocket.
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// ErrorMessage is the object that is returned when there is an error via
// websocket.
type ErrorMessage struct {
	Msg      string `json:"msg"`
	HookName string `json:"hookname"`
}

// FindHandler is a function used to find a specific handler in the router
type FindHandler func(string) (Handler, bool)

// Client is the object to handle the websocket connection.
type Client struct {
	send         chan Message
	socket       *websocket.Conn
	findHandler  FindHandler
	stopChannels map[int]chan bool
	Pocket       *pocket.Pocket
	Code         string
	AccessToken  string
}

// NewStopChannel creates a new stop channel for a given websocket conncection
func (c *Client) NewStopChannel(stopKey int) chan bool {
	c.StopForKey(stopKey)
	stop := make(chan bool)
	c.stopChannels[stopKey] = stop
	return stop
}

// StopForKey stops a socket conncection for given name
func (c *Client) StopForKey(key int) {
	if ch, found := c.stopChannels[key]; found {
		ch <- true
		delete(c.stopChannels, key)
	}
}

// Read reads any messages sent to our client
func (c *Client) Read() {
	var message Message
	for {
		_, r, err := c.socket.NextReader()
		if err != nil {
			msg := fmt.Sprintf("Error getting reader: %s", err.Error())
			clientLog.Criticalf(msg)
			// c.SendError("system", errors.New(msg))
			break
		}

		d := json.NewDecoder(r)
		d.UseNumber()

		if err := d.Decode(&message); err != nil {
			msg := fmt.Sprintf("Error decoding message: %s", err.Error())
			clientLog.Criticalf(msg)
			c.SendError("system", errors.New(msg))
			break
		}

		if fn, ok := c.findHandler(message.Name); ok {
			clientLog.Infof("Received request for %s\n", message.Name)
			fn(c, message.Data)
		}
	}

	c.socket.Close()
}

// Write sends a message back to any websocket that is listening
func (c *Client) Write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			clientLog.Criticalf("Error writing JSON: %s", err.Error())
			continue
		}
	}

	c.socket.Close()
}

// Close sends a message to close all open channels.
func (c *Client) Close() {
	for _, ch := range c.stopChannels {
		ch <- true
	}
	close(c.send)
}

// SendError sends a error message back to whomever is listening
func (c *Client) SendError(name string, err error) {
	e := ErrorMessage{
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
