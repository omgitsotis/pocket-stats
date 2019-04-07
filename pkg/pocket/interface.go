package pocket

// PocketClient is the interface to talk to the Pocket client
type PocketClient interface {
	GetAuth(uri string) (string, error)
	ReceieveAuth(key string) (*User, error)
	IsAuthed() bool
	GetArticles(since int) (RetrieveResult, error)
	DebugGetArticles(since int) ([]byte, error)
}
