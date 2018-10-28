package model

// AuthRequest is the body of the request used to get the auth from the pocket
type AuthLinkRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

// RequestToken is the body of the response of the get auth request
type RequestToken struct {
	Code string `json:"code"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Link is the response body which holds the link used to authorise a user
type Link struct {
	URL string `json:"url"`
}

// AuthoriseRequest is the body of the request to auth the user.
type AuthRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}
