package dao

import "github.com/omgitsotis/pocket-stats/pocket/model"

type DAO interface {
	AddUser(string) (int64, error)
	AddArticle(model.Row) error
	IsUser(int64) (bool, error)
	GetCountForDates(int64, int64) (*model.Stats, error)
	GetArticle(int64) (*model.Row, error)
	UpdateArticle(*model.Row) error
	GetLastAdded() (int64, error)
}
