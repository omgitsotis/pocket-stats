package pocket

type Request struct {
    ConsumerKey string `json:"consumer_key"`
    RedirectURI string `json:"redirect_uri"`
}

type RequestToken struct {
    Code string `json:"code"`
}
