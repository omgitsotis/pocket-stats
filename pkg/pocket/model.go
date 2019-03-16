package pocket

// ItemStatus represents the read status of the item.
type ItemStatus int

const (
	ItemStatusUnread   ItemStatus = 0
	ItemStatusArchived ItemStatus = 1
	ItemStatusDeleted  ItemStatus = 2
)

type RetrieveResult struct {
	List     map[string]Article
	Status   int
	Complete int
	Since    int
}

// GetArticleList converts the map in RetrieveResult to a list
func (r *RetrieveResult) GetArticleList() []Article {
	aList := make([]Article, 0)

	for _, value := range r.List {
		aList = append(aList, value)
	}

	return aList
}

// Article represents the article object from the pocket API
type Article struct {
	ItemID        int        `json:"item_id,string"`
	ResolvedId    int        `json:"resolved_id,string"`
	ResolvedURL   string     `json:"resolved_url"`
	GivenTitle    string     `json:"given_title"`
	ResolvedTitle string     `json:"resolved_title"`
	Favorite      int        `json:",string"`
	Status        ItemStatus `json:",string"`
	WordCount     int        `json:"word_count,string"`
	Tags          map[string]map[string]interface{}
	Authors       map[string]map[string]interface{}
	SortId        int    `json:"sort_id"`
	TimeAdded     string `json:"time_added"`
	TimeUpdated   string `json:"time_updated"`
	TimeRead      string `json:"time_read"`
}

type Sort string

const (
	SortNewest Sort = "newest"
	SortOldest Sort = "oldest"
	SortTitle  Sort = "title"
	SortSite   Sort = "site"
)

// RetrieveOption is the options for retrieve API.
type RetrieveOption struct {
	State       string `json:"state,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Sort        Sort   `json:"sort,omitempty"`
	DetailType  string `json:"detailType,omitempty"`
	Since       int    `json:"since,omitempty"`
	Count       int    `json:"count,omitempty"`
	Offset      int    `json:"offset,omitempty"`
	AccessToken string `json:"access_token"`
	ConsumerKey string `json:"consumer_key"`
}
