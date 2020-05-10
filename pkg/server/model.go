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

// WordsPerMinute is the number of words per minute I can currently read
const WordsPerMinute = 146

// Stats is the response to the GetStats call
type Stats struct {
	Totals        *StatTotals    `json:"totals"`
	Itemised      *ItemisedStats `json:"itemised"`
	Tags          *TagStats      `json:"tags"`
	PreviousStats *PreviousStats `json:previous`
}

// StatTotals returns the totals of the articles updated within the time range
type StatTotals struct {
	ArticlesRead  int64 `json:"articles_read"`
	ArticlesAdded int64 `json:"articles_added"`
	WordsRead     int64 `json:"words_read"`
	WordsAdded    int64 `json:"words_added"`
	TimeRead      int64 `json:"time_read"`
	TimeAdded     int64 `json:"time_added"`
}

type TagTotals struct {
	ArticlesRead int64 `json:"articles_read"`
	WordsRead    int64 `json:"words_read"`
	TimeRead     int64 `json:"time_read"`
}

type PreviousStats struct {
	Totals *StatTotals `json:"totals"`
	Tags   *TagStats   `json:"tags"`
}

// ItemisedStats is a map of the day to the totals of the articles updated
type ItemisedStats map[int64]*StatTotals

// TagStats is a map of the tags to the totals of the articles updated
type TagStats map[string]*TagTotals
