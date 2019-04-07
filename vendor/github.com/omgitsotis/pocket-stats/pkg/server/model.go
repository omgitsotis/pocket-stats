package server

type updateResponse struct {
	Date int64 `json:"date_updated"`
}

type healthcheckResp struct {
	Status string `json:"status"`
}

// Link is the response body which holds the link used to authorise a user
type Link struct {
	URL string `json:"url"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
