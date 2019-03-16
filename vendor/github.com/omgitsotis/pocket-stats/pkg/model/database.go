package model

type HealthcheckResp struct {
	Status string `json:"status"`
}

// User represents the user in the database
type User struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	ID          int64  `json:"id"`
	LastUpdated int64  `json:"last_updated"`
}

type Article struct {
	Name     string
	DateRead int64
	Tags     []string
}
