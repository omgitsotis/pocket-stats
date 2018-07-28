package model

const Added = "0"
const Archived = "1"
const Deleted = "2"

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
	ID          int64  `json:"id"`
	LastUpdated int64  `json:"last_updated"`
}

type DataParam struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
	Since       int64  `json:"since"`
	State       string `json:"state"`
	Sort        string `json:"sort"`
	Type        string `json:"detailType"`
}

// DataList is the list of articles retrieved from the Pocket API
type DataList struct {
	Status int             `json:"status"`
	Values map[string]Data `json:"list"`
}

// Data is each pocket item saved
type Data struct {
	ItemID     string         `json:"item_id"`
	ResolvedID string         `json:"resolved_id"`
	Status     string         `json:"status"`
	WordCount  string         `json:"word_count"`
	Title      string         `json:"given_title"`
	Tags       map[string]Tag `json:"tags"`
}

type Tag struct {
	ItemID string `json:"item_id"`
	Tag    string `json:"tag"`
}

type Article struct {
	ID        int64
	DateAdded int64
	DateRead  int64
	WordCount int64
	Status    string
	UserID    int64
	Tag       string
}

type InputParams struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
	Date  int64  `json:"date"`
}

type CountRow struct {
	Date      int64 `json:"date"`
	WordCount int64 `json:"word_count"`
}

type Stats struct {
	Start      int64            `json:"start_date"`
	End        int64            `json:"end_date"`
	DateValues map[int64]*Stat  `json:"date_values"`
	TagValues  map[string]*Stat `json:"tag_values"`
	Totals     TotalStats       `json:"totals"`
}

// Stat is the struct that holds the statistics for either a date or a tag
type Stat struct {
	ArticleAdded int64 `json:"articles_added"`
	ArticleRead  int64 `json:"articles_read"`
	WordsAdded   int64 `json:"words_added"`
	WordsRead    int64 `json:"words_read"`
	TimeReading  int64 `json:"time_reading"`
}

type TotalStats struct {
	ArticlesAdded int64 `json:"total_articles_added"`
	ArticlesRead  int64 `json:"total_articles_read"`
	WordsAdded    int64 `json:"total_words_added"`
	WordsRead     int64 `json:"total_words_read"`
	TimeReading   int64 `json:"total_time_reading"`
}

type StatsParams struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
