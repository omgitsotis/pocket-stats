package dao

import "github.com/omgitsotis/pocket-stats/pocket/model"

type DAO interface {
	AddUser(string) (int64, error)
	AddArticle(model.Article) error
	IsUser(int64) (bool, error)
	GetArticles(int64, int64) ([]model.Article, error)
	GetArticle(int64) (*model.Article, error)
	UpdateArticle(*model.Article) error
	GetLastAdded() (int64, error)
}
