package pocket

// RequestToken is the body of the response of the get auth request
type RequestToken struct {
	Code string `json:"code"`
}

// AuthRequest is the body of the request to auth the user.
type AuthRequest struct {
	ConsumerKey string `json:"consumer_key"`
	Code        string `json:"code"`
}

// AuthLinkRequest is the body of the request used to get the auth from the pocket
type AuthLinkRequest struct {
	ConsumerKey string `json:"consumer_key"`
	RedirectURI string `json:"redirect_uri"`
}

// User is the user represented in the Pocket API
type User struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	ID          int64  `json:"id"`
	LastUpdated int64  `json:"last_updated"`
}

// ItemStatus represents the read status of the item.
type ItemStatus int

const (
	// ItemStatusUnread is an item marked as unread in the Pocket API
	ItemStatusUnread ItemStatus = 0
	// ItemStatusArchived is an item marked as Archived in the Pocket API
	ItemStatusArchived ItemStatus = 1
	// ItemStatusDeleted is an item marked as Deleted in the Pocket API
	ItemStatusDeleted ItemStatus = 2
)

// RetrieveResult is the response body of the Get request to the Pocket API
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
		if value.Status != ItemStatusDeleted {
			aList = append(aList, value)
		}
	}

	return aList
}

// Article represents the article object from the pocket API
type Article struct {
	ItemID        int        `json:"item_id,string"`
	ResolvedID    int        `json:"resolved_id,string"`
	ResolvedURL   string     `json:"resolved_url"`
	GivenTitle    string     `json:"given_title"`
	ResolvedTitle string     `json:"resolved_title"`
	Favorite      int        `json:",string"`
	Status        ItemStatus `json:",string"`
	WordCount     int        `json:"word_count,string"`
	Tags          map[string]map[string]interface{}
	Authors       map[string]map[string]interface{}
	SortID        int    `json:"sort_id"`
	TimeAdded     string `json:"time_added"`
	TimeUpdated   string `json:"time_updated"`
	TimeRead      string `json:"time_read"`
}

// Sort represents the ways you can sort in the Pocket API
type Sort string

const (
	// SortNewest is the string used to sort items by newest added
	SortNewest Sort = "newest"
	// SortOldest is the string used to sort itmes by oldest added
	SortOldest Sort = "oldest"
	// SortTitle is the string used to sort items by title
	SortTitle Sort = "title"
	// SortSite is the string used to sort itmes by site
	SortSite Sort = "site"
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
