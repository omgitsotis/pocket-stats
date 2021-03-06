package database

import "github.com/omgitsotis/pocket-stats/pkg/pocket"

// DBCLient is the interface to the database
type DBCLient interface {
	SaveArticles(articles []pocket.Article) error
	GetArticle(id int) (*Article, error)
	GetLastUpdateDate() (int, error)
	UpsertArticles(articles []pocket.Article) error
	SaveUpdateDate(int64) error
	GetArticlesByDate(start, end int64) ([]Article, error)
	GetArticlesByTag(start, end int64, tag string) ([]Article, error)
	DeleteArticle(int) error
}
