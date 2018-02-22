package pocket

type Request struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

type RequestToken struct {
	Code string `json:"code"`
}

type Authorise struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

type User struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}

type DataParam struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Since       int64  `json:"since"`
	State       string `json:"state"`
	Sort        string `json:"sort"`
}

type DataList struct {
	Status int             `json:"status"`
	Values map[string]Data `json:"list"`
}

type Data struct {
	ItemID     string `json:"item_id"`
	ResolvedID string `json:"resolved_id"`
	Status     string `json:"status"`
	WordCount  string `json:"word_count"`
}
