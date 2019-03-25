package database

import "github.com/omgitsotis/pocket-stats/pkg/pocket"

// DBCLient is the interface to the database
type DBCLient interface {
	SaveArticles(articles []pocket.Article) error
	GetArticle(id string) (Article, error)
}
