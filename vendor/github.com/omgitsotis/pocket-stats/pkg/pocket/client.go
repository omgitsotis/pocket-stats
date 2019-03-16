package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/omgitsotis/pocket-stats/pkg/model"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

const pocketURL = "https://getpocket.com/v3"

// The only user that is allowed to use this app. (Me, nwah nwah nwah)
const Username = "omgitsotis"

func Init(l *logrus.Logger) {
	log = l
}

type Client struct {
	consumerID string
	client     *http.Client
	authedUser *model.User
}

func New(consumerID string, cli *http.Client) *Client {
	return &Client{
		consumerID: consumerID,
		client:     cli,
	}
}

// GetAuth gets the request token from pocket
func (c *Client) GetAuth(uri string) (string, error) {
	r := model.AuthLinkRequest{c.consumerID, uri}
	var rt model.RequestToken
	if err := c.call("/oauth/request", r, &rt); err != nil {
		return "", err
	}

	log.Debugf("repsone code returned [%s]", rt.Code)
	return rt.Code, nil
}

// RecievedAuth gets the access token from pocket, and returns the user from the
// database
func (c *Client) ReceieveAuth(key string) (*model.User, error) {
	a := model.AuthRequest{c.consumerID, key}
	var user model.User

	if err := c.call("/oauth/authorize", a, &user); err != nil {
		return nil, err
	}

	if user.Username != Username {
		return nil, errors.New("Unauthorised user for this app")
	}

	c.authedUser = &user

	// TODO move this logic out of here
	// date, err := p.dao.GetLastAdded()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// user.LastUpdated = date
	// logger.Infof("Last added date: [%d]", date)

	return &user, nil
}

// IsAppAuthed returns a boolean based on whether the client has an authorised
// user or not
func (c *Client) IsAuthed() bool {
	return c.authedUser != nil
}

func (c *Client) GetArticles(offset int) (*RetrieveResult, error) {
	req := RetrieveOption{
		Count:       100,
		Sort:        SortOldest,
		DetailType:  "complete",
		ContentType: "article",
		State:       "all",
		AccessToken: c.authedUser.AccessToken,
		ConsumerKey: c.consumerID,
		Offset:      offset,
		Since:       1422403200,
	}

	log.Debugf("Params to send %+v", req)

	var resp RetrieveResult

	if err := c.call("/get", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// call makes api requests to the Pocket api and marshal the results.
func (c *Client) call(uri string, body, t interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		log.Errorf("Error marshalling params: %s", err.Error())
		return err
	}

	uri = fmt.Sprintf("%s%s", pocketURL, uri)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Errorf("Error creating request: %s", err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		log.Errorf("error performing request: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Errorf("Status [%s] Error %s", res.Status, res.Header["X-Error"])
		return errors.New(res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(t)
	if err != nil {
		log.Errorf("Error decoding body: %s", err.Error())
		return err
	}

	return nil
}
