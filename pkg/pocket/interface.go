package pocket

import "github.com/omgitsotis/pocket-stats/pkg/model"

type PocketClient interface {
	GetAuth(uri string) (string, error)
	ReceieveAuth(key string) (*model.User, error)
	IsAuthed() bool
	GetArticles(offset int) (RetrieveResult, error)
}
