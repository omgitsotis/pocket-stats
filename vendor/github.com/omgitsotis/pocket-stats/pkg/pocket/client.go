package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/omgitsotis/pocket-stats/pkg/model"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

const pocketURL = "https://getpocket.com/v3"

// Username is the only user that is allowed to use this app. (Me, nwah nwah nwah)
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

// ReceieveAuth gets the access token from pocket, and returns the user from the
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

	return &user, nil
}

// IsAppAuthed returns a boolean based on whether the client has an authorised
// user or not
func (c *Client) IsAuthed() bool {
	return c.authedUser != nil
}

func (c *Client) GetArticles(since int) (*RetrieveResult, error) {
	req := RetrieveOption{
		Sort:        SortOldest,
		DetailType:  "complete",
		ContentType: "article",
		State:       "all",
		AccessToken: c.authedUser.AccessToken,
		ConsumerKey: c.consumerID,
		Since:       since,
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

	url := fmt.Sprintf("%s%s", pocketURL, uri)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "error reading body")
	}

	err = json.Unmarshal(data, t)
	if err != nil {
		if uri == "/get" {
			log.Warn("Got unmarshal error, trying with empty list")
			t, err = tryEmptyResponse(data)
			if err != nil {
				return errors.Wrap(err, "error decoding body")
			}
		}

		return errors.Wrap(err, "error decoding body")
	}

	return nil
}

func tryEmptyResponse(data []byte) (interface{}, error) {
	type emptyResponse struct {
		List     []Article
		Status   int
		Complete int
		Since    int
	}

	var er emptyResponse

	if err := json.Unmarshal(data, &er); err != nil {
		return nil, err
	}

	return &RetrieveResult{
		Status:   er.Status,
		Complete: er.Status,
		Since:    er.Since,
		List:     map[string]Article{},
	}, nil
}
