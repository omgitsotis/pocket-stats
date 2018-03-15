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

type DataList struct {
	Status int             `json:"status"`
	Values map[string]Data `json:"list"`
}

type Data struct {
	ItemID     string `json:"item_id"`
	ResolvedID string `json:"resolved_id"`
	Status     string `json:"status"`
	WordCount  string `json:"word_count"`
	Title      string `json:"given_title"`
}

type Article struct {
	ID        int64
	DateAdded int64
	DateRead  int64
	WordCount int64
	Status    string
	UserID    int64
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
	Start  int64           `json:"start_date"`
	End    int64           `json:"end_date"`
	Value  map[int64]*Stat `json:"value"`
	Totals TotalStats      `json:"totals"`
}

type Stat struct {
	ArticleAdded int64 `json:"articles_added"`
	ArticleRead  int64 `json:"articles_read"`
	WordAdded    int64 `json:"words_added"`
	WordRead     int64 `json:"words_read"`
}

type TotalStats struct {
	ArticlesAdded int64 `json:"total_articles_added"`
	ArticlesRead  int64 `json:"total_articles_read"`
	WordsAdded    int64 `json:"total_words_added"`
	WordsRead     int64 `json:"total_words_read"`
}

type StatsParams struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
