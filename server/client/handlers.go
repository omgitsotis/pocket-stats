package client

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/omgitsotis/pocket-stats/server/pocket/model"
)

// sendAuth gets a request token from Pocket and returns the link required to
// auth the user.
func sendAuth(client *Client, data interface{}) {
	code, err := client.Pocket.GetAuth("http://localhost:4000/auth/recieved")
	if err != nil {
		client.SendError("auth link", err)
		return
	}

	client.Code = code

	u := fmt.Sprintf(
		"https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
		code,
		"http://localhost:4000/auth/recieved",
	)

	type Link struct {
		URL string `json:"url"`
	}

	link := Link{u}
	client.send <- Message{"auth link", link}
}

// loadUser gets the user from the data. This method is called when the access
// token is found in the cookies.
func loadUser(client *Client, data interface{}) {
	type AccessToken struct {
		Token string `json:"token"`
	}

	var token AccessToken
	err := mapstructure.Decode(data, &token)
	if err != nil {
		client.SendError("auth cached", err)
		return
	}

	user, err := client.Pocket.GetUser()
	if err != nil {
		client.SendError("auth cached", err)
		return
	}

	clientLog.Debugf("Received token [%s]", token)
	user.AccessToken = token.Token

	client.send <- Message{"auth cached", user}
}

// initDB loads the database with article information from pocket from a given
// date.
func initDB(client *Client, data interface{}) {
	var params model.InputParams
	err := mapstructure.Decode(data, &params)
	if err != nil {
		clientLog.Errorf("Error decoding params: %s", err.Error())
		client.SendError("data init", err)
		return
	}

	err = client.Pocket.InitDB(params)
	if err != nil {
		client.SendError("auth init", err)
		return
	}

	client.send <- Message{"data init", "Complete"}
}

// getStatistics gets the article statistics between two given dates
func getStatistics(client *Client, data interface{}) {
	var p model.StatsParams

	err := mapstructure.Decode(data, &p)
	if err != nil {
		clientLog.Errorf("Error decoding params: %s", err.Error())
		client.SendError("data get", err)
		return
	}

	clientLog.Infof("Get stats from %d to %d", p.Start, p.End)

	stats, err := client.Pocket.GetStatsForDates(p)
	if err != nil {
		client.SendError("data get", err)
		return
	}

	client.send <- Message{"data get", stats}
}

// updateDB will load article information from pocket from the last update date
func updateDB(client *Client, data interface{}) {
	var params model.InputParams
	err := mapstructure.Decode(data, &params)
	if err != nil {
		clientLog.Errorf("Error decoding params: %s", err.Error())
		client.SendError("data update", err)
		return
	}

	date, err := client.Pocket.UpdateDB(params)
	if err != nil {
		client.SendError("data update", err)
		return
	}

	client.send <- Message{"data update", date}
}

// loadData loads the last week's worth of statistics
func loadData(client *Client, data interface{}) {
	stats, err := client.Pocket.LoadData()
	if err != nil {
		client.SendError("data get", err)
		return
	}

	client.send <- Message{"data get", stats}
}
